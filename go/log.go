package hyperkit

import (
	golog "log"
)

// Logger is an interface for logging.
type Logger interface {
	// Debugf logs a message with "debug" severity (very verbose).
	Debugf(format string, v ...interface{})
	// Infof logs a message with "info" severity (less verbose).
	Infof(format string, v ...interface{})
	// Fatalf logs a fatal error message, and exits 1.
	Fatalf(format string, v ...interface{})
}

// StandardLogger makes the go standard logger comply to our Logger interface.
type StandardLogger struct{}

// Debugf logs a message with "debug" severity.
func (*StandardLogger) Debugf(f string, v ...interface{}) {
	golog.Printf(f, v...)
}

// Infof logs a message with "info" severity.
func (*StandardLogger) Infof(f string, v ...interface{}) {
	golog.Printf(f, v...)
}

// Fatalf logs a fatal error message, and exits 1.
func (*StandardLogger) Fatalf(f string, v ...interface{}) {
	golog.Fatalf(f, v...)
}

// Log receives stdout/stderr of the hyperkit process itself, if set.
// It defaults to the go standard logger.
var log Logger = &StandardLogger{}

// SetLogger sets the logger to use.
func SetLogger(l Logger) {
	log = l
}
