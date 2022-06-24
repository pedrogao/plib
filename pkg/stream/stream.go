package stream

import (
	"sort"
	"sync"

	"github.com/pedrogao/plib/pkg/collection"
	"github.com/pedrogao/plib/pkg/common"
	"github.com/pedrogao/plib/pkg/task"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

type (
	// Item of stream
	Item interface {
		comparable
	}

	// Stream computing
	Stream[T Item] struct {
		source <-chan T // 只读 channel，不能写
	}

	// KeyFunc key生成器
	// item - stream中的元素
	KeyFunc[T Item, R Item] func(item T) R

	// FilterFunc 过滤函数
	FilterFunc[T Item] func(item T) bool

	// MapFunc 对象转换函数
	MapFunc[T Item, R Item] func(item T) R

	// LessFunc 对象比较
	LessFunc[T Item] func(a, b T) bool

	// WalkFunc 遍历函数
	WalkFunc[T Item, R Item] func(item T, pip chan<- R)

	// PredicateFunc 匹配函数
	PredicateFunc[T Item] func(item T) bool

	// ForAllFunc 对所有元素执行操作
	ForAllFunc[T Item] func(pip <-chan T)

	// ForEachFunc 对每个item执行操作
	ForEachFunc[T Item] func(item T)

	// GenerateFunc item生成函数
	GenerateFunc[T Item] func(source <-chan T)
)

// Distinct 去重，使用 map 来实现去重
func (s Stream[T]) Distinct(keyFunc KeyFunc[T, T]) Stream[T] {
	source := make(chan T)
	common.GoSafe(func() { // 新建协程写数据
		// channel记得关闭是个好习惯
		defer close(source)
		keys := make(map[T]common.PlaceholderType)

		for item := range s.source {
			// 自定义去重逻辑
			key := keyFunc(item) // 这里的key类型是R
			// 如果key不存在,则将数据写入新的channel
			if _, ok := keys[key]; !ok {
				source <- item
				keys[key] = common.Placeholder
			}
		}
	})
	return Range[T](source)
}

func (s Stream[T]) Filter(filterFunc FilterFunc[T], opts ...Option) Stream[T] {
	return s.Walk(func(item T, pip chan<- T) {
		if filterFunc(item) {
			pip <- item
		}
	}, opts...)
}

func (s Stream[T]) Walk(fn WalkFunc[T, T], opts ...Option) Stream[T] {
	option := buildOptions(opts...)
	if option.unlimitedWorkers {
		return s.walkUnLimited(fn, option)
	}
	return s.walkLimited(fn, option)
}

func (s Stream[T]) walkUnLimited(fn WalkFunc[T, T],
	option *rxOptions) Stream[T] {
	// 创建带缓冲区的channel
	// 默认为16,channel中元素超过16将会被阻塞
	pipe := make(chan T, defaultWorkers)
	go func() {
		var wg sync.WaitGroup

		for item := range s.source {
			// 需要读取s.source的所有元素
			// 这里也说明了为什么channel最后写完记得完毕
			// 如果不关闭可能导致协程一直阻塞导致泄漏
			// 重要, 不赋值给val是个典型的并发陷阱，后面在另一个goroutine里使用了
			val := item
			wg.Add(1)
			// 安全模式下执行函数
			common.GoSafe(func() {
				defer wg.Done()
				fn(val, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	}()

	// 返回新的Stream
	return Range[T](pipe)
}

func (s Stream[T]) walkLimited(fn WalkFunc[T, T],
	option *rxOptions) Stream[T] {
	pipe := make(chan T, option.workers)
	go func() {
		var wg sync.WaitGroup
		// 控制协程数量
		pool := make(chan common.PlaceholderType, option.workers)

		for item := range s.source {
			// 重要, 不赋值给val是个典型的并发陷阱，后面在另一个goroutine里使用了
			val := item
			// 超过协程限制时将会被阻塞
			pool <- common.Placeholder
			// 这里也说明了为什么channel最后写完记得完毕
			// 如果不关闭可能导致协程一直阻塞导致泄漏
			wg.Add(1)

			// 安全模式下执行函数
			common.GoSafe(func() {
				defer func() {
					wg.Done()
					//执行完成后读取一次pool释放一个协程位置
					<-pool
				}()
				fn(val, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	}()
	return Range[T](pipe)
}

func (s Stream[T]) Head(n int64) Stream[T] {
	if n < 1 {
		panic("n must be greater than 1")
	}
	source := make(chan T)
	go func() {
		for item := range s.source {
			n--
			// n值可能大于s.source长度,需要判断是否>=0
			if n >= 0 {
				source <- item
			}
			// let successive method go ASAP even we have more items to skip
			// why we don't just break the loop, because if break,
			// this former goroutine will block forever, which will cause goroutine leak.
			// n==0说明source已经写满可以进行关闭了
			// 既然source已经满足条件了为什么不直接进行break跳出循环呢?
			// 作者提到了防止协程泄漏
			// 因为每次操作最终都会产生一个新的Stream,旧的Stream永远也不会被调用了
			if n == 0 {
				close(source)
				break
			}
		}
		// 上面的循环跳出来了说明n大于s.source实际长度
		// 依旧需要显示关闭新的source
		if n > 0 {
			close(source)
		}
	}()
	return Range[T](source)
}

func (s Stream[T]) Tail(n int64) Stream[T] {
	if n < 1 {
		panic("n must be greater than 1")
	}
	source := make(chan T)
	go func() {
		ring := collection.NewRing[T](int(n))
		// 读取全部元素，如果数量>n环形切片能实现新数据覆盖旧数据
		// 保证获取到的一定最后n个元素
		for item := range s.source {
			ring.Add(item)
		}
		for _, item := range ring.Take() {
			source <- item
		}
		close(source)
	}()
	return Range[T](source)
}

func (s Stream[T]) Map(fn MapFunc[T, T], opts ...Option) Stream[T] {
	return s.Walk(func(item T, pip chan<- T) {
		pip <- fn(item)
	}, opts...)
}

func (s Stream[T]) Reverse() Stream[T] {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}
	for i := len(items)/2 - 1; i >= 0; i-- {
		opp := len(items) - 1 - i
		items[i], items[opp] = items[opp], items[i]
	}
	return Just[T](items...)
}

func (s Stream[T]) Sort(fn LessFunc[T]) Stream[T] {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return fn(items[i], items[j])
	})
	return Just[T](items...)
}

func (s Stream[T]) AllMatch(fn PredicateFunc[T]) bool {
	for item := range s.source {
		if !fn(item) {
			// 需要排空 s.source，否则前面的goroutine可能阻塞
			go drain(s.source)
			return false
		}
	}

	return true
}

func (s Stream[T]) AnyMatch(fn PredicateFunc[T]) bool {
	for item := range s.source {
		if fn(item) {
			// 需要排空 s.source，否则前面的goroutine可能阻塞
			go drain(s.source)
			return true
		}
	}

	return false
}

func (s Stream[T]) NoneMatch(fn func(item T) bool) bool {
	for item := range s.source {
		if fn(item) {
			// 需要排空 s.source，否则前面的goroutine可能阻塞
			go drain(s.source)
			return false
		}
	}

	return true
}

func (s Stream[T]) Count() int {
	var count int
	for range s.source {
		count++
	}
	return count
}

func (s Stream[T]) Done() {
	// 排空 channel，防止 goroutine 阻塞泄露
	drain(s.source)
}

func (s Stream[T]) ForAll(fn ForAllFunc[T]) {
	fn(s.source)
}

func (s Stream[T]) ForEach(fn ForEachFunc[T]) {
	for item := range s.source {
		fn(item)
	}
}

func (s Stream[T]) channel() <-chan T {
	return s.source
}

func Range[T Item](source <-chan T) Stream[T] {
	return Stream[T]{
		source: source,
	}
}

func Just[T Item](items ...T) Stream[T] {
	source := make(chan T, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)
	return Range[T](source)
}

func From[T Item](generate GenerateFunc[T]) Stream[T] {
	source := make(chan T)
	common.GoSafe(func() {
		defer close(source)
		generate(source)
	})
	return Range[T](source)
}

func (s Stream[T]) Concat(steams ...Stream[T]) Stream[T] {
	// 创建新的无缓冲channel
	source := make(chan T)
	go func() {
		// 创建一个waiGroup对象
		group := task.NewRoutineGroup()
		// 异步从原channel读取数据
		group.Run(func() {
			for item := range s.source {
				source <- item
			}
		})
		// 异步读取待拼接Stream的channel数据
		for _, stream := range steams {
			// 每个Stream开启一个协程
			group.Run(func() {
				for item := range stream.channel() {
					source <- item
				}
			})
		}
		// 阻塞等待读取完成
		group.Wait()
		close(source)
	}()
	// 返回新的Stream
	return Range[T](source)
}

// drain drains the given channel.
func drain[T any](channel <-chan T) {
	for range channel { // 消费 channel
	}
}
