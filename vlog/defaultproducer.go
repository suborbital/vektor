package vlog

import (
	"fmt"
	"strings"
)

type defaultProducer struct{}

// ErrorString prints a string as an error
func (d *defaultProducer) ErrorString(msgs ...string) string {
	return fmt.Sprintf("(E) %s", strings.Join(msgs, " "))
}

// Error prints a string as an error
func (d *defaultProducer) Error(err error) string {
	return fmt.Sprintf("(E) %s", err.Error())
}

// Warn prints a string as an warning
func (d *defaultProducer) Warn(msgs ...string) string {
	return fmt.Sprintf("(W) %s", strings.Join(msgs, " "))
}

// Info prints a string as an info message
func (d *defaultProducer) Info(msgs ...string) string {
	return fmt.Sprintf("(I) %s", strings.Join(msgs, " "))
}

// Debug prints a string as debug output
func (d *defaultProducer) Debug(msgs ...string) string {
	return fmt.Sprintf("(D) %s", strings.Join(msgs, " "))
}

// Trace prints a function name and returns a function to be deferred, logging the completion of a function
func (d *defaultProducer) Trace(fnName string) (string, func() string) {
	traceFunc := func() string {
		return (fmt.Sprintf("(T) %s completed", fnName))
	}

	return (fmt.Sprintf("(T) %s", fnName)), traceFunc
}
