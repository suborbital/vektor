package vk

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/suborbital/vektor/vlog"
)

// Response represents a non-error HTTP response
type Response struct {
	status int
	body   interface{}
}

// Respond returns a filled-in response
func Respond(status int, body interface{}) Response {
	r := Response{
		status: status,
		body:   body,
	}

	return r
}

// R is `Respond` for those who prefer terse code
func R(status int, body interface{}) Response {
	return Respond(status, body)
}

// TODO: add convenience helpers for status codes

const (
	contentTypeJSON        contentType = "application/json"
	contentTypeTextPlain   contentType = "text/plain"
	contentTypeOctetStream contentType = "application/octet-stream"
)

// converts _something_ into bytes, best it can:
// if data is Response type, returns (status, body processed as below)
// if bytes, return (200, bytes)
// if string, return (200, []byte(string))
// if struct, return (200, json(struct))
// otherwise, return (500, nil)
func responseOrOtherToBytes(l *vlog.Logger, data interface{}) (int, []byte, contentType) {
	if data == nil {
		return http.StatusNoContent, []byte{}, contentTypeTextPlain
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
		return statusCode, b, contentTypeOctetStream
	} else if s, ok := realData.(string); ok {
		return statusCode, []byte(s), contentTypeTextPlain
	}

	// otherwise, assume it's a struct of some kind,
	// so JSON marshal it and return it
	json, err := json.Marshal(realData)
	if err != nil {
		l.Error(errors.Wrap(err, "failed to Marshal response struct"))

		return genericErrorResponseCode, []byte(genericErrorResponseBytes), contentTypeTextPlain
	}

	return statusCode, json, contentTypeJSON
}

type RawString string
type RawBytes []byte

// RespondWeb converts a value to either raw string or json, and sends it to the client. Ctx is a placeholder here, it
// is currently unused, but will be used for tracing / logging help purposes later.
func RespondWeb(_ context.Context, w http.ResponseWriter, data any, statusCode int) error {
	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// If the data is deliberately set as raw bytes or raw strings, return them as is. If we don't have this, they will
	// be marshalled into json, which will wrap a naked string into double quotes.
	switch i := data.(type) {
	case RawBytes:
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		_, _ = w.Write(i)
		return nil
	case RawString:
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(i))
		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
