package vk

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Response represents a non-error HTTP response
type Response struct {
	status int
	body   interface{}
}

// Reply returns a filled-in response
func Reply(status int, body interface{}) Response {
	r := Response{
		status: status,
		body:   body,
	}

	return r
}

// R is `Reply` for those who prefer terse code
func R(status int, body interface{}) Response {
	return Reply(status, body)
}

// TODO: add convenience helpers for status codes

// converts _something_ into bytes, best it can:
// if data is Response type, returns (status, body processed as below)
// if bytes, return (200, bytes)
// if string, return (200, []byte(string))
// if struct, return (200, json(struct))
// otherwise, return (500, nil)
func responseOrOtherToBytes(data interface{}) (int, []byte) {
	if data == nil {
		return http.StatusNoContent, []byte{}
	}

	statusCode := http.StatusOK
	realData := data

	// first, check if it's response type, and unpack it for further processing
	if r, ok := data.(Response); ok {
		statusCode = r.status
		realData = r.body
	}

	// if data is []byte or string, return it as-is
	if b, ok := realData.([]byte); ok {
		return statusCode, b
	} else if s, ok := realData.(string); ok {
		return statusCode, []byte(s)
	}

	// otherwise, assume it's a struct of some kind,
	// so JSON marshal it and return it
	json, err := json.Marshal(realData)
	if err != nil {
		return 500, []byte(errors.Wrap(err, "failed to json Marshal response struct").Error()) // TODO: make this error reporting better
	}

	return statusCode, json
}
