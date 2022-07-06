package vk

import (
	"fmt"
)

// Error is an interface representing a failed request
type Error interface {
	Error() string // this ensures all Errors will also conform to the normal error interface

	Message() string
	Status() int
}

// ErrorResponse is a concrete implementation of Error,
// representing a failed HTTP request
type ErrorResponse struct {
	StatusCode  int    `json:"status"`
	MessageText string `json:"message"`
}

// Error returns a full error string
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.MessageText)
}

// Status returns the error status code
func (e *ErrorResponse) Status() int {
	return e.StatusCode
}

// Message returns the error's message
func (e *ErrorResponse) Message() string {
	return e.MessageText
}

// Err returns an error with status and message
func Err(status int, message string) Error {
	e := &ErrorResponse{
		StatusCode:  status,
		MessageText: message,
	}

	return e
}

// E is Err for those who like terse code
func E(status int, message string) Error {
	return Err(status, message)
}

// Wrap wraps an error in vk.Error
func Wrap(status int, err error) Error {
	return Err(status, err.Error())
}
