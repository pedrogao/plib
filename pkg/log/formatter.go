package log

// Formatter which format log items and output
type Formatter interface {
	Format(e *Entry) error
}
