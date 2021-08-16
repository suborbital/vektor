package vtest

import (
	"bytes"
	"net/http"
)

type fakeResponseWriter struct {
	header http.Header
	*bytes.Buffer
	status int
}

// newFakeResponseWriter creates a FakeResponseWriter for capturing Vektor responses
func newFakeResponseWriter() *fakeResponseWriter {
	return &fakeResponseWriter{make(http.Header), &bytes.Buffer{}, 0}
}

// Header implements http.ResponseWriter
func (w *fakeResponseWriter) Header() http.Header {
	return w.header
}

// WriteHeader implements http.ResponseWriter
func (w *fakeResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w fakeResponseWriter) Status() int {
	return w.status
}
