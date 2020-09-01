package vlog

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Producer represents an object that is considered a producer of messages
type Producer interface {
	ErrorString(...string) string         // Logs an error string
	Error(error) string                   // Logs an error obj
	Warn(...string) string                // Logs a warning
	Info(...string) string                // Logs information
	Debug(...string) string               // Logs debug information
	Trace(string) (string, func() string) // Logs a function call and returns a function to be deferred, indicating the end of the function
}

// Logger is the main logger object, responsible for taking input from the
// producer and managing scoped loggers
type Logger struct {
	producer Producer
	opts     Options
	output   io.Writer
}

// Default returns a Logger using the default producer
func Default(opts ...OptionsModifier) *Logger {
	prod := &defaultProducer{}

	return New(prod, opts...)
}

// New returns a Logger with the provided producer and options
func New(producer Producer, opts ...OptionsModifier) *Logger {
	options := newOptions(opts...)

	v := &Logger{
		producer: producer,
		opts:     options,
	}

	output, err := outputForOptions(options)
	if err != nil {
		v.output = os.Stdout
		os.Stderr.Write([]byte("failed to set vlog output: " + err.Error() + "\n"))
	} else {
		v.output = output
	}

	return v
}

// ErrorString logs a string as an error
func (v *Logger) ErrorString(msgs ...string) {
	msg := v.producer.ErrorString(msgs...)

	v.log(msg, nil, 1)
}

// Error logs an error as an error
func (v *Logger) Error(err error) {
	msg := v.producer.Error(err)

	v.log(msg, nil, 1)
}

// Warn logs a string as an warning
func (v *Logger) Warn(msgs ...string) {
	msg := v.producer.Warn(msgs...)

	v.log(msg, nil, 2)
}

// Info logs a string as an info message
func (v *Logger) Info(msgs ...string) {
	msg := v.producer.Info(msgs...)

	v.log(msg, nil, 3)
}

// Debug logs a string as debug output
func (v *Logger) Debug(msgs ...string) {
	msg := v.producer.Debug(msgs...)

	v.log(msg, nil, 4)
}

// Trace logs a function name and returns a function to be deferred, logging the completion of a function
func (v *Logger) Trace(fnName string) func() {
	msg, traceFunc := v.producer.Trace(fnName)

	v.log(msg, nil, 5)

	return func() {
		msg := traceFunc()

		v.log(msg, nil, 5)
	}
}

func (v *Logger) log(message string, scope interface{}, level int) {
	if level > v.opts.level {
		return
	}

	// send the raw message to the console
	if v.output != os.Stdout {
		// throwing away the error here since there's nothing much we can do
		os.Stdout.Write([]byte(message))
		os.Stdout.Write([]byte("\n"))
	}

	structured := structuredLog{
		LogMessage: message,
		Timestamp:  time.Now(),
		Level:      level,
		AppMeta:    v.opts.appMeta,
		ScopeMeta:  scope,
	}

	structuredJSON, err := json.Marshal(structured)
	if err != nil {
		os.Stderr.Write([]byte("[vlog] failed to marshal structured log"))
	}

	_, err = v.output.Write(structuredJSON)
	if err != nil {
		os.Stderr.Write([]byte("failed to write to configured output: " + err.Error() + "\n"))
	} else {
		v.output.Write([]byte("\n"))
	}

}

func outputForOptions(opts Options) (io.Writer, error) {
	var output io.Writer

	if opts.filepath != "" {
		file, err := os.OpenFile(opts.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}

		output = file
	} else {
		output = os.Stdout
	}

	return output, nil
}
