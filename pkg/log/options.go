package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

// FmtEmptySeparate log format empty separator
const FmtEmptySeparate = ""

// Level for log
type Level uint8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel // panic but recover
	FatalLevel // process will crash
)

var LevelNameMapping = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	PanicLevel: "PANIC",
	FatalLevel: "FATAL",
}

var errUnmarshalNilLevel = errors.New("can't unmarshal a nil Level")

func (l *Level) unmarshalText(text []byte) bool {
	switch string(text) {
	case "debug", "DEBUG":
		*l = DebugLevel
	case "info", "INFO", "": // make the zero value useful
		*l = InfoLevel
	case "warn", "WARN":
		*l = WarnLevel
	case "error", "ERROR":
		*l = ErrorLevel
	case "panic", "PANIC":
		*l = PanicLevel
	case "fatal", "FATAL":
		*l = FatalLevel
	default:
		return false
	}
	return true
}

// UnmarshalText unmarshal text to a level.
func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized level: %q", text)
	}
	return nil
}

type options struct {
	output        io.Writer
	level         Level
	stdLevel      Level
	formatter     Formatter
	disableCaller bool
}

// Option for logger
type Option func(*options)

func initOptions(opts ...Option) (o *options) {
	o = &options{}

	for _, opt := range opts {
		opt(o)
	}

	if o.output == nil {
		o.output = os.Stderr
	}

	if o.formatter == nil {
		o.formatter = &TextFormatter{}
	}

	return
}

func WithLevel(level Level) Option {
	return func(o *options) {
		o.level = level
	}
}

func WithOutput(output io.Writer) Option {
	return func(o *options) {
		o.output = output
	}
}

func WithFormatter(formatter Formatter) Option {
	return func(o *options) {
		o.formatter = formatter
	}
}

func WithDisableCaller(disableCaller bool) Option {
	return func(o *options) {
		o.disableCaller = disableCaller
	}
}
