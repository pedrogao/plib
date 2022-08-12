package log

import (
	"io"
	"os"
	"sync"
	"unsafe"
)

var std = New()

type logger struct {
	mu        sync.Mutex
	entryPool *sync.Pool
	opts      *options
}

func New(opts ...Option) *logger {
	logger := &logger{opts: initOptions(opts...)}
	logger.entryPool = &sync.Pool{
		// if entry is not enough, then call entry
		New: func() any {
			return entry(logger)
		},
	}
	return logger
}

func SetOptions(opts ...Option) {
	std.SetOptions(opts...)
}

func (l *logger) SetOptions(opts ...Option) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, opt := range opts {
		opt(l.opts)
	}
}

func (l *logger) entry() *Entry {
	return l.entryPool.Get().(*Entry)
}

func Writer() io.Writer {
	return std
}

func (l *logger) Write(data []byte) (n int, err error) {
	l.entry().write(l.opts.stdLevel, FmtEmptySeparate, *(*string)(unsafe.Pointer(&data)))
	return 0, nil
}

func (l *logger) Debug(args ...any) {
	l.entry().write(DebugLevel, FmtEmptySeparate, args...)
}

func (l *logger) Info(args ...any) {
	l.entry().write(InfoLevel, FmtEmptySeparate, args...)
}

func (l *logger) Warn(args ...any) {
	l.entry().write(WarnLevel, FmtEmptySeparate, args...)
}

func (l *logger) Error(args ...any) {
	l.entry().write(ErrorLevel, FmtEmptySeparate, args...)
}

func (l *logger) Panic(args ...any) {
	l.entry().write(PanicLevel, FmtEmptySeparate, args...)
}

func (l *logger) Fatal(args ...any) {
	l.entry().write(FatalLevel, FmtEmptySeparate, args...)
	os.Exit(1)
}

func (l *logger) Debugf(format string, args ...any) {
	l.entry().write(DebugLevel, format, args...)
}

func (l *logger) Infof(format string, args ...any) {
	l.entry().write(InfoLevel, format, args...)
}

func (l *logger) Warnf(format string, args ...any) {
	l.entry().write(WarnLevel, format, args...)
}

func (l *logger) Errorf(format string, args ...any) {
	l.entry().write(ErrorLevel, format, args...)
}

func (l *logger) Panicf(format string, args ...any) {
	l.entry().write(PanicLevel, format, args...)
}

func (l *logger) Fatalf(format string, args ...any) {
	l.entry().write(FatalLevel, format, args...)
	os.Exit(1)
}

// std logger

func Debug(args ...any) {
	std.Debug(args...)
}

func Info(args ...any) {
	std.Info(args...)
}

func Warn(args ...any) {
	std.Warn(args...)
}

func Error(args ...any) {
	std.Error(args...)
}

func Panic(args ...any) {
	std.Error(args...)
}

func Fatal(args ...any) {
	std.Fatal(args...)
}

func Debugf(format string, args ...any) {
	std.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	std.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	std.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	std.Errorf(format, args...)
}

func Panicf(format string, args ...any) {
	std.Panicf(format, args...)
}

func Fatalf(format string, args ...any) {
	std.Fatalf(format, args...)
}
