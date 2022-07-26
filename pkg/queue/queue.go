package queue

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"syscall"
	"time"
)

const (
	itemLengthSize = 4
)

var (
	Full  = errors.New("queue is full")
	Empty = errors.New("queue is empty")
)

type (
	Queue interface {
		// Push data to queue
		Push(string) error

		// Pop data from queue
		Pop() (string, error)

		// Empty check queue is empty
		Empty() bool
	}

	metadata struct {
		size      int     // 队列大小
		chunkSize int     // 块大小
		head      *cursor // 头部块
		tail      *cursor // 尾部块
	}

	cursor struct {
		num    int // 文件序号
		offset int // 文件 offset
		length int // 文件 length
	}

	queueOptions struct {
		MaxSize   int
		ChunkSize int
		TempDir   string
		AutoSave  bool
		GCTimeout time.Duration
	}

	Option func(*queueOptions)

	diskQueue struct {
		path       string
		maxSize    int
		chunkSize  int
		tempDir    string
		autoSave   bool
		gcTimeout  time.Duration
		meta       *metadata
		serializer Serializer
		headFile   *os.File
		tailFile   *os.File
		gcTicker   *time.Ticker
	}
)

func WithChunkSize(chunkSize int) Option {
	return func(ops *queueOptions) {
		ops.ChunkSize = chunkSize
	}
}

func WithTempDir(tempDir string) Option {
	return func(ops *queueOptions) {
		ops.TempDir = tempDir
	}
}

func WithMaxSize(maxSize int) Option {
	return func(ops *queueOptions) {
		ops.MaxSize = maxSize
	}
}

func WithAutoSave(autoSave bool) Option {
	return func(ops *queueOptions) {
		ops.AutoSave = autoSave
	}
}

func WithGCTimeout(timeout time.Duration) Option {
	return func(ops *queueOptions) {
		ops.GCTimeout = timeout
	}
}

var defaultOptions = func() queueOptions {
	return queueOptions{
		MaxSize:   0,
		ChunkSize: 100, // TODO
		AutoSave:  false,
	}
}

func New(path string, options ...Option) (*diskQueue, error) {
	if path == "" {
		return nil, fmt.Errorf("path: %s not valid", path)
	}

	ops := defaultOptions()
	for _, option := range options {
		option(&ops)
	}
	if ops.TempDir == "" {
		ops.TempDir = os.TempDir()
	}

	q := &diskQueue{
		path:       path,
		maxSize:    ops.MaxSize,
		chunkSize:  ops.ChunkSize,
		tempDir:    ops.TempDir,
		autoSave:   ops.AutoSave,
		serializer: NewJsonSerializer(),
		gcTimeout:  ops.GCTimeout,
		gcTicker:   time.NewTicker(ops.GCTimeout),
	}

	err := q.init()
	if err != nil {
		return nil, err
	}

	go q.gc()

	return q, nil
}

func (q *diskQueue) Push(val string) error {
	if q.qSize() >= q.maxSize {
		return Full
	}
	data := []byte(val)
	// 写到 tail 文件
	if q.meta.tail.length+len(data) > q.meta.chunkSize {
		if err := q.advanceTail(); err != nil {
			return fmt.Errorf("advance tail file err: %s", err)
		}
	}
	lbuf := make([]byte, itemLengthSize+len(data))
	binary.BigEndian.PutUint32(lbuf, uint32(len(data)))
	copy(lbuf[itemLengthSize:], data)
	// TODO handle n != len(data) ?
	/*n,*/ _, err := q.tailFile.Write(lbuf)
	if err != nil {
		return fmt.Errorf("write queue data err: %s", err)
	}
	q.meta.tail.length += len(data)
	q.meta.tail.offset += 1
	q.meta.size++

	// fsync
	err = syscall.Fsync(int(q.tailFile.Fd()))
	if err != nil {
		return fmt.Errorf("flush queue data file err: %s", err)
	}

	return nil
}

func (q *diskQueue) advanceTail() error {
	// 1. 将旧的 tail file 刷盘
	var err error
	err = syscall.Fsync(int(q.tailFile.Fd()))
	if err != nil {
		return fmt.Errorf("flush queue data file err: %s", err)
	}
	// 2. 创建新的 tail file
	tailPath := q.qFile(q.meta.tail.num + 1)
	q.tailFile, err = os.OpenFile(tailPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("open tail file err: %s", err)
	}
	// 3. TODO 保存 旧 tail file 的 meta
	// 4. 更新新的 tail meta
	q.meta.tail.num += 1
	q.meta.tail.length = 0
	q.meta.tail.offset = 0
	return nil
}

func (q *diskQueue) Pop() (string, error) {
	// 检查 tail 是否超过了 head
	err := q.checkEmpty()
	if err != nil {
		return "", err
	}
	if q.meta.head.length >= q.meta.chunkSize {
		if err := q.advanceHead(); err != nil {
			return "", fmt.Errorf("advance head file err: %s", err)
		}
	}

	lbuf := make([]byte, itemLengthSize)
	// TODO handle n != len(data) ?
	_, err = q.headFile.Read(lbuf)
	if err != nil {
		return "", fmt.Errorf("read queue file err: %s", err)
	}
	length := binary.BigEndian.Uint32(lbuf)
	dbuf := make([]byte, length)
	// TODO handle n != len(data) ?
	_, err = q.headFile.Read(dbuf)
	if err != nil {
		return "", fmt.Errorf("read queue file err: %s", err)
	}
	q.meta.head.offset -= 1
	q.meta.head.length -= itemLengthSize + int(length)
	q.meta.size--
	return string(dbuf), nil
}

func (q *diskQueue) checkEmpty() error {
	if q.meta.tail.num > q.meta.head.num {
		return Empty
	}
	if q.meta.tail.num == q.meta.head.num &&
		q.meta.tail.offset >= q.meta.tail.offset {
		return Empty
	}
	return nil
}

func (q *diskQueue) advanceHead() error {
	// 1. 将旧的 head file 刷盘
	var err error
	err = syscall.Fsync(int(q.headFile.Fd()))
	if err != nil {
		return fmt.Errorf("flush queue data file err: %s", err)
	}
	// 2. 创建新的 head file
	headPath := q.qFile(q.meta.head.num + 1)
	q.headFile, err = os.OpenFile(headPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("open head file err: %s", err)
	}
	// 3. TODO 保存 旧 head file 的 meta
	// 4. 更新新的 head meta
	q.meta.head.num += 1
	q.meta.head.length = 0
	q.meta.head.offset = 0
	return nil
}

func (q *diskQueue) Empty() bool {
	return q.qSize() == 0
}

func (q *diskQueue) init() error {
	if _, err := os.Stat(q.path); os.IsNotExist(err) {
		_ = os.MkdirAll(q.path, os.ModePerm)
	}
	if _, err := os.Stat(q.tempDir); os.IsNotExist(err) {
		_ = os.MkdirAll(q.path, os.ModePerm)
	}
	// load meta data
	q.meta = q.loadMeta()

	var err error
	headPath := q.qFile(q.meta.head.num)
	q.headFile, err = os.OpenFile(headPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("open head file err: %s", err)
	}

	tailPath := q.qFile(q.meta.tail.num)
	q.tailFile, err = os.OpenFile(tailPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("open tail file err: %s", err)
	}

	return nil
}

func (q *diskQueue) qFile(number int) string {
	return path.Join(q.path, fmt.Sprintf("q%05d", number))
}

//
// metadata save & load
//

func (q *diskQueue) saveMeta() error {
	temp, err := os.CreateTemp(q.tempDir, "meta")
	if err != nil {
		return fmt.Errorf("create meta temp file err: %s", err)
	}

	err = q.serializer.DumpFile(temp, q.meta)
	if err != nil {
		return fmt.Errorf("dump meta to file err: %s", err)
	}

	p := q.metaPath()
	err = os.Rename(temp.Name(), p) // FIXME?
	if err != nil {
		return fmt.Errorf("repalce meta file err: %s", err)
	}

	return nil
}

func (q *diskQueue) loadMeta() *metadata {
	p := q.metaPath()
	meta := &metadata{
		size:      0,
		chunkSize: q.chunkSize,
		head: &cursor{
			num:    0,
			offset: 0,
			length: 0,
		},
		tail: &cursor{
			num:    0,
			offset: 0,
			length: 0,
		},
	}
	if ok := exists(p); ok {
		if err := q.serializer.Load(p, meta); err == nil {
			return meta
		}
	}
	return meta
}

func (q *diskQueue) metaPath() string {
	return path.Join(q.path, "meta")
}

func (q *diskQueue) qSize() int {
	return q.meta.size
}

func (q *diskQueue) gc() {
	ticker := q.gcTicker
	for {
		select {
		case <-ticker.C:
			q.cleanFiles()
		}
	}
}

func (q *diskQueue) cleanFiles() {
	if q.meta.head.num == 0 {
		return
	}
	// chunk size 的作用就在这里，清理的时候，可以容忍一定长度上的浪费
	for i := 0; i < q.meta.head.num; i++ {
		abandonPath := q.qFile(i)
		if exist := exists(abandonPath); exist {
			err := os.Remove(abandonPath)
			if err != nil {
				log.Printf("[gc] remove unsed file err: %s", err)
			}
		}
	}
	return
}
