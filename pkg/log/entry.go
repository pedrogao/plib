package log

import (
	"bytes"
	"runtime"
	"strings"
	"time"
)

// Entry store log information
type Entry struct {
	logger *logger
	Buffer *bytes.Buffer
	Map    map[string]any
	Level  Level
	Time   time.Time
	File   string
	Line   int
	Func   string
	Format string
	Args   []any
}

func entry(logger *logger) *Entry {
	return &Entry{
		logger: logger,
		Buffer: new(bytes.Buffer),
		Map:    make(map[string]any, 5),
	}
}

func (e *Entry) write(level Level, format string, args ...any) {
	if e.logger.opts.level > level {
		return
	}
	e.Level = level
	e.Time = time.Now()
	e.Format = format
	e.Args = args
	if !e.logger.opts.disableCaller {
		if pc, file, line, ok := runtime.Caller(3); !ok {
			e.File = "???"
			e.Func = "???"
		} else {
			e.File, e.Line, e.Func = file, line, runtime.FuncForPC(pc).Name()
			e.Func = e.Func[strings.LastIndex(e.Func, "/")+1:]
		}
	}
	e.format()
	e.writer()
	e.release()
}

func (e *Entry) format() {
	_ = e.logger.opts.formatter.Format(e)
}

func (e *Entry) writer() {
	e.logger.mu.Lock()
	defer e.logger.mu.Unlock()

	_, _ = e.logger.opts.output.Write(e.Buffer.Bytes())
}

func (e *Entry) release() {
	e.Args, e.Line, e.File, e.Format, e.Func = nil, 0, "", "", ""
	e.Buffer.Reset()
	e.logger.entryPool.Put(e)
}
