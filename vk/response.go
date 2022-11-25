package vk

import (
	"encoding/json"
	"net/http"
)

// RespondJSON converts a value to json, and sends it to the client. Ctx is a placeholder here, it is currently unused,
// but will be used for tracing / logging help purposes later.
func RespondJSON(w http.ResponseWriter, data any, statusCode int) error {
	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	return respondBytes(w, jsonData, statusCode)
}

// JSON is a shorthand for RespondJSON for people who like terse code
func JSON(w http.ResponseWriter, data any, statusCode int) error {
	return RespondJSON(w, data, statusCode)
}

// RespondString sends the data as a raw string to the client. Ctx is a placeholder here, it is currently unused, but
// will be used for tracing / logging help purposes later.
func RespondString(w http.ResponseWriter, data string, statusCode int) error {
	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	w.Header().Set("Content-Type", "text/plain")

	return respondBytes(w, []byte(data), statusCode)
}

// RespondBytes takes the content we want to send back to the client as byte slice, does an early exit for no content,
// and then straight pipes into the private respondBytes. This is in front of the private func because of pattern
// consistency and making sure that a 204 with not empty content does not get written to the ResponseWriter.
func RespondBytes(w http.ResponseWriter, data []byte, statusCode int) error {
	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	w.Header().Set("Content-Type", "text/plain")

	return respondBytes(w, data, statusCode)
}

// respondBytes takes content as []byte and pipes it into the http.ResponseWriter. The assumption is that the content
// type header has already been set. The only time this is not called is when the status code is http.StatusNoContent,
// which would have broken early in the callers.
func respondBytes(w http.ResponseWriter, content []byte, statusCode int) error {
	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(content); err != nil {
		return err
	}

	return nil
}
