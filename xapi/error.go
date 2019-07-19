package x

import (
	"encoding/json"
	"fmt"
)

// Error represents an HTTP error
type Error struct {
	Status  int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

// Err returns an error with status and message
func Err(status int, message string) Error {
	e := Error{
		Status:  status,
		Message: message,
	}

	return e
}

// E is Err for those who like terse code
func E(status int, message string) Error {
	return Err(status, message)
}

// converts _something_ into bytes, best it can:
// if data is Error type, returns (status, {status: status, message: message})
// if other error, returns (500, []byte(err.Error()))
func errorOrOtherToBytes(err error) (int, []byte) {
	statusCode := 500
	realData := []byte(err.Error())

	// first, check if it's response type, and unpack it for further processing
	if e, ok := err.(Error); ok {
		statusCode = e.Status
		errJSON, marshalErr := json.Marshal(e)
		if marshalErr != nil {
			return statusCode, realData
		}

		return statusCode, errJSON
	}

	return statusCode, realData
}
