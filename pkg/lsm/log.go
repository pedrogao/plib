package lsm

import (
	"fmt"
	"os"
)

// AppendLog 追加日志
type AppendLog struct {
	filename string
	stream   *os.File
}

// NewAppendLog 新建追加日志
func NewAppendLog(filename string) (*AppendLog, error) {
	// link: https://en.wikipedia.org/wiki/File-system_permissions
	stream, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("open log file: %s err: %s", filename, err)
	}
	log := &AppendLog{
		filename: filename,
		stream:   stream,
	}
	return log, nil
}

func (l *AppendLog) WriteString(val string) error {
	_, err := l.stream.WriteString(val)
	if err != nil {
		return fmt.Errorf("write log file err: %s", err)
	}

	err = l.stream.Sync()
	if err != nil {
		return fmt.Errorf("flush log file err: %s", err)
	}

	return nil
}

func (l *AppendLog) Write(val []byte) error {
	_, err := l.stream.Write(val)
	if err != nil {
		return fmt.Errorf("write log file err: %s", err)
	}

	err = l.stream.Sync()
	if err != nil {
		return fmt.Errorf("flush log file err: %s", err)
	}

	return nil
}

func (l *AppendLog) Clear() error {
	err := l.stream.Close()
	if err != nil {
		return fmt.Errorf("close log file err: %s", err)
	}

	l.stream, err = os.OpenFile(l.filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("reopen log file err: %s", err)
	}

	return nil
}
