package queue

import (
	"fmt"
	"os"
	"path"
)

type (
	Queue any

	metadata struct {
		size      int     // 未完成任务数
		chunkSize int     // 块大小
		head      *cursor // 头部块
		tail      *cursor // 尾部块
	}

	cursor struct {
		num, count, offset int
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
		ChunkSize: 100,
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
	q := &diskQueue{
		path:       path,
		maxSize:    ops.MaxSize,
		chunkSize:  ops.ChunkSize,
		tempDir:    ops.TempDir,
		autoSave:   ops.AutoSave,
		serializer: NewJsonSerializer(),
	}

	q.init()

	return q, nil
}

func (q *diskQueue) init() {
	if _, err := os.Stat(q.path); os.IsNotExist(err) {
		_ = os.MkdirAll(q.path, os.ModePerm)
	}
}

func (q *diskQueue) qFile(number int) string {
	return path.Join(q.path, fmt.Sprintf("q%05d", number))
}

//
// metadata save & load
//

func (q *diskQueue) saveMeta() error {
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
