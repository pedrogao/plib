package workpool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gammazero/deque"
	"github.com/pedrogao/plib/pkg/log"
)

/**
reference:
- http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang
- http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html
*/

const (
	// 休眠时间，如果 worker 长时间不工作，那么停掉 worker
	// If workers idle for at least this period of time, then stop a worker.
	idleTimeout = 2 * time.Second
)

type (
	// Task interface for all
	Task interface {
		// Do 执行任务
		Do() error

		// Context 返回上下文
		Context() context.Context
	}

	// defaultTask 默认任务
	defaultTask struct {
		fn  func()
		ctx context.Context
	}
)

// Do 执行任务
func (d *defaultTask) Do() error {
	d.fn()
	return nil
}

// Context 返回上下文
func (d *defaultTask) Context() context.Context {
	return d.ctx
}

// New an instance of worker pool
func New(maxWorkers int) *WorkerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	pool := &WorkerPool{
		maxWorkers:  maxWorkers,
		taskQueue:   make(chan Task, 1),
		workerQueue: make(chan Task),
		stopSignal:  make(chan struct{}),
		stoppedChan: make(chan struct{}),
	}

	go pool.dispatch()

	return pool
}

// WorkerPool a collection of concurrent goroutine workers
// 任务执行流：
// +-----------+     +--------------+     +-------------+
// | taskQueue | --> | waitingQueue | --> | workerQueue |
// +-----------+     +--------------+     +-------------+
type WorkerPool struct {
	maxWorkers   int           // 最大 worker 个数
	taskQueue    chan Task     // 任务队列
	workerQueue  chan Task     // 工作队列
	stoppedChan  chan struct{} // 已停止状态
	stopSignal   chan struct{} // 停止动作
	waitingQueue deque.Deque   // 等待队列
	stopLock     sync.Mutex
	stopOnce     sync.Once
	stopped      bool  // 是否停止
	waiting      int32 // 等待任务个数
	wait         bool  // 是否等待
}

// Size worker 个数
func (p *WorkerPool) Size() int {
	return p.maxWorkers
}

// Stop 停止且无等待
func (p *WorkerPool) Stop() {
	p.stop(false)
}

// StopWait 停止且等待任务完成
func (p *WorkerPool) StopWait() {
	p.stop(true)
}

// Stopped 是否已停止
func (p *WorkerPool) Stopped() bool {
	p.stopLock.Lock()
	defer p.stopLock.Unlock()

	return p.stopped
}

// Submit 提交任务无阻塞
func (p *WorkerPool) Submit(task Task) {
	if task == nil {
		return
	}
	p.taskQueue <- task // 无需等待
}

// SubmitWait 提交任务且已阻塞
func (p *WorkerPool) SubmitWait(task Task) {
	if task == nil {
		return
	}
	// 通过 channel 来阻塞等待任务完成
	doneChan := make(chan struct{})
	t := &defaultTask{
		fn: func() {
			err := task.Do()
			if err != nil {
				log.Error("do task err: ", err)
			}
			close(doneChan) // 任务执行完毕后通知 doneChan，然后函数返回
		},
		ctx: task.Context(),
	}
	p.taskQueue <- t
	<-doneChan // 阻塞，等待被关闭
}

// Pause 暂停所有 worker
// 思路：向所有 worker 提交一个等待任务，等待时间由 ctx 设置
// 如果 ctx 为 emptyCtx 那么将直接退出
func (p *WorkerPool) Pause(ctx context.Context) {
	p.stopLock.Lock()
	defer p.stopLock.Unlock()
	if p.stopped {
		return
	}

	ready := new(sync.WaitGroup)
	ready.Add(p.maxWorkers)
	for i := 0; i < p.maxWorkers; i++ {
		t := &defaultTask{
			fn: func() {
				ready.Done()
				// 使用 select 阻塞 worker，直到 ctx 撤销、或者 stop
				select {
				case <-ctx.Done():
				case <-p.stopSignal:
				}
			},
			ctx: ctx,
		}
		p.Submit(t)
	}
	ready.Wait() // wait for all
}

// WaitingQueueSize 等待队列任务个数
func (p *WorkerPool) WaitingQueueSize() int {
	return int(atomic.LoadInt32(&p.waiting))
}

func (p *WorkerPool) dispatch() {
	defer close(p.stoppedChan) // 关闭 stopped channel，通知 stop 函数退出
	timeout := time.NewTimer(idleTimeout)
	var (
		workerCount int
		idle        bool
		wg          sync.WaitGroup
	)

Loop:
	for {
		// 如果等待队列有任务，那么先将等待队列中的任务转交给工作队列中
		if p.waitingQueue.Len() != 0 {
			if !p.processWaitingQueue() {
				break Loop // 如果 taskQueue 被关闭了，那么直接退出loop
			}
			// 每次任务转移，那么 continue，继续下次
			// 这里相当于 while 循环，直到将 waitingQueue 中的任务清零
			continue
		}
		// 此处如果 waitingQueue 队列中的任务清零后，那么尝试将 taskQueue 中的
		// 任务提交至 workerQueue
		select {
		case task, ok := <-p.taskQueue: // 从任务队列中拿到任务
			if !ok {
				break Loop // 任务队列已关闭，那么直接 break
			}
			// 将任务提交至 workerQueue
			select {
			case p.workerQueue <- task: // 提交至工作队列
			default:
				// 如果 workerQueue 提交不了，那么尝试新建 worker
				// 或者将任务提交至 waitingQueue
				if workerCount < p.maxWorkers {
					wg.Add(1)
					go startWorker(task, p.workerQueue, &wg)
					workerCount++
				} else {
					// 如果 worker 已经达到了最大数，那么将任务交给等待队列
					p.waitingQueue.PushBack(task)
					atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len()))
				}
			}
			idle = false // 无空闲 worker
		case <-timeout.C:
			// taskQueue 中长时间无任务，那么尝试杀死多余的 worker，避免浪费
			if idle && workerCount > 0 {
				// 向 workerQueue 中发送 nil，杀死 worker
				if p.killIdleWorker() {
					workerCount--
				}
			}
			idle = true                // 有空闲 worker
			timeout.Reset(idleTimeout) // 重置定时器
		}
	}
	// 等待处理 waitingQueue 中的任务
	if p.wait {
		p.runQueuedTasks()
	}
	// 向任务队列中发送 nil 任务，worker 收到 nil 后会将自己杀死
	for workerCount > 0 {
		// 为什么可以严格按照 nil 的个数来杀死所有的 worker？
		// 因此 worker 严格按照收到 nil 后直接退出，那么一个 nil 就能杀死一个 worker
		// 这样就不会有 goroutine 泄漏的问题
		p.workerQueue <- nil
		workerCount--
	}
	wg.Wait() // 等待 worker 死亡，避免 goroutine 泄漏
	timeout.Stop()
}

// startWorker 开启 worker
func startWorker(task Task, workerQueue chan Task, wg *sync.WaitGroup) {
	err := task.Do() // 1. 先执行提交任务
	if err != nil {
		log.Error("do task err: ", err)
	}
	// TODO go safety
	go worker(workerQueue, wg) // 2. 然后开启循环，监听 workerQueue
}

func worker(workerQueue chan Task, wg *sync.WaitGroup) {
	for task := range workerQueue {
		// 收到 nil，退出函数，即将自己杀死
		if task == nil {
			wg.Done() // -1
			return
		}
		err := task.Do() // 否则，执行任务
		if err != nil {
			log.Error("do task err: ", err)
		}
	}
}

func (p *WorkerPool) stop(wait bool) {
	// 关闭，有且执行一次
	p.stopOnce.Do(func() {
		// 关闭 chan，会通知所有 reader、writer
		// 通知 pool 已经关闭
		close(p.stopSignal)
		p.stopLock.Lock()
		p.stopped = true
		p.stopLock.Unlock()
		p.wait = wait // 设置是否等待
		// 关闭任务队列后后不会再接收任务
		close(p.taskQueue)
	})
	// 阻塞 stop 函数，直到 dispatch 函数完成
	<-p.stoppedChan
}

func (p *WorkerPool) processWaitingQueue() bool {
	select {
	case task, ok := <-p.taskQueue: // 将任务从 taskQueue 中推出
		if !ok { // task queue is closed, so return false
			return false
		}
		p.waitingQueue.PushBack(task) // 然后加入到等待队列
	case p.workerQueue <- p.waitingQueue.Front().(Task): // 或者将任务从等待队列中推出，然后加入工作队列
		// 顶部 job pop
		p.waitingQueue.PopFront()
	}
	atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len())) // 更新等待队列个数
	return true
}

func (p *WorkerPool) killIdleWorker() bool {
	select {
	// send nil to worker queue
	case p.workerQueue <- nil:
		return true
	default:
		return false
	}
}

func (p *WorkerPool) runQueuedTasks() {
	for p.waitingQueue.Len() != 0 {
		// 等待队列任务中的job出队，然后加入到 worker 队列
		p.workerQueue <- p.waitingQueue.PopFront().(Task)
		atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len()))
	}
}
