package glog

import (
	"fmt"
	"strings"
)

// Logger represents an object that is considered a logger
type Logger interface {
	ErrorString(...string) // Logs an error string
	Error(error)           // Logs an error obj
	Warn(...string)        // Logs a warning
	Info(...string)        // Logs information
	Debug(...string)       // Logs debug information
	Trace(string) func()   // Logs a function call and returns a function to be deferred, indicating the end of the function
	// Sensitive(string)      // Logs sensitive information, should be used wisely
}

// LogLevelTrace and others represent log levels
const (
	LogLevelTrace = "trace" // 5
	LogLevelDebug = "debug" // 4
	LogLevelInfo  = "info"  // 3
	LogLevelWarn  = "warn"  // 2
	LogLevelError = "error" // 1
)

// GLoggerOptions represents the options for a GLogger
type GLoggerOptions struct {
	level    int
	filepath string
	prefix   string
}

// GLogger is the "builtin" implementation of a GusT Logger
type GLogger struct {
	opts *GLoggerOptions
}

// DefaultLogger returns a GLogger that prints to the console
func DefaultLogger() *GLogger {
	return LoggerWithOptions(defaultOptions())
}

// LoggerWithOptions returns a logger based on the provided options
func LoggerWithOptions(opts *GLoggerOptions) *GLogger {
	l := &GLogger{
		opts: opts,
	}

	return l
}

// ErrorString prints a string as an error
func (gl *GLogger) ErrorString(msgs ...string) {
	gl.log(fmt.Sprintf("(E) %s", strings.Join(msgs, " ")))
}

// Error prints a string as an error
func (gl *GLogger) Error(err error) {
	gl.log(fmt.Sprintf("(E) %s", err.Error()))
}

// Warn prints a string as an warning
func (gl *GLogger) Warn(msgs ...string) {
	if gl.opts.level >= 2 {
		gl.log(fmt.Sprintf("(W) %s", strings.Join(msgs, " ")))
	}
}

// Info prints a string as an info message
func (gl *GLogger) Info(msgs ...string) {
	if gl.opts.level >= 3 {
		gl.log(fmt.Sprintf("(I) %s", strings.Join(msgs, " ")))
	}
}

// Debug prints a string as debug output
func (gl *GLogger) Debug(msgs ...string) {
	if gl.opts.level >= 4 {
		gl.log(fmt.Sprintf("(D) %s", strings.Join(msgs, " ")))
	}
}

// Trace prints a function name and returns a function to be deferred, logging the completion of a function
func (gl *GLogger) Trace(fnName string) func() {
	if gl.opts.level >= 5 {
		gl.log(fmt.Sprintf("(T) %s", fnName))

		return func() {
			gl.log(fmt.Sprintf("(T) %s completed", fnName))
		}
	}

	return func() {}
}

func (gl *GLogger) log(msg string) {
	if gl.isFileLogger() {
		// TODO: add file logging
	}

	if gl.opts.prefix != "" {
		fmt.Printf("%s %s\n", gl.opts.prefix, msg)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

func (gl *GLogger) isFileLogger() bool {
	return gl.opts.filepath != ""
}

func defaultOptions() *GLoggerOptions {
	o := &GLoggerOptions{
		level:    logLevelValFromString(LogLevelInfo),
		filepath: "",
		prefix:   "",
	}

	return o
}

func logLevelValFromString(level string) int {
	switch level {
	case LogLevelTrace:
		return 5
	case LogLevelDebug:
		return 4
	case LogLevelInfo:
		return 3
	case LogLevelWarn:
		return 2
	case LogLevelError:
		return 1
	}

	return 3
}
