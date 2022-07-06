package vk

import (
	"context"
	"encoding/json"
	"net/http"
)

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
