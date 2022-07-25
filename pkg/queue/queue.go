package queue

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"syscall"
)

const (
	itemLengthSize = 4
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
		count  int // 消息个数
		offset int // 文件 offset
	}

	queueOptions struct {
		MaxSize   int
		ChunkSize int
		TempDir   string
		AutoSave  bool
	}

	Option func(*queueOptions)

	diskQueue struct {
		path       string
		maxSize    int
		chunkSize  int
		tempDir    string
		autoSave   bool
		meta       *metadata
		serializer Serializer
		headFile   *os.File
		tailFile   *os.File
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
	}

	err := q.init()
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (q *diskQueue) Push(val string) error {
	data := []byte(val)
	// 写到 tail 文件
	if q.meta.tail.offset+len(data) > q.meta.chunkSize {
		q.meta.tail.num += 1
		// TODO reopen new tail file
	}
	// TODO 编码 |len|data...|
	lbuf := make([]byte, itemLengthSize+len(data))
	binary.BigEndian.PutUint32(lbuf, uint32(len(data)))
	copy(lbuf[itemLengthSize:], data)
	// TODO handle n != len(data) ?
	/*n,*/ _, err := q.tailFile.Write(lbuf)
	if err != nil {
		return fmt.Errorf("write queue data err: %s", err)
	}
	q.meta.tail.offset += len(data)
	q.meta.tail.count += 1

	// fsync
	err = syscall.Fsync(int(q.tailFile.Fd()))
	if err != nil {
		return fmt.Errorf("flush queue data file err: %s", err)
	}

	return nil
}

func (q *diskQueue) Pop() (string, error) {
	if q.meta.head.offset >= q.meta.chunkSize {
		q.meta.head.num += 1
		// TODO reopen head file
	}

	lbuf := make([]byte, itemLengthSize)
	// TODO handle n != len(data) ?
	_, err := q.headFile.Read(lbuf)
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
	q.meta.head.count -= 1
	q.meta.head.offset -= itemLengthSize + int(length)
	// TODO gc
	return string(dbuf), nil
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

	// TODO load file & meta if need
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
			count:  0,
			offset: 0,
		},
		tail: &cursor{
			num:    0,
			count:  0,
			offset: 0,
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
