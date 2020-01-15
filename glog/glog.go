package glog

// Logger represents an object that is considered a logger
type Logger interface {
	ErrorString(...string) // Logs an error string
	Error(error)           // Logs an error obj
	Warn(...string)        // Logs a warning
	Info(...string)        // Logs information
	Debug(...string)       // Logs debug information
	Trace(string) func()   // Logs a function call and returns a function to be deferred, indicating the end of the function
	Sensitive(string)      // Logs sensitive information, should be used wisely
}
