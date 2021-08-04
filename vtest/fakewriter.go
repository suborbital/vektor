package vtest

import (
	"bytes"
	"net/http"
)

type FakeResponseWriter struct {
	header http.Header
	*bytes.Buffer
	status int
}

// NewFakeResponseWriter creates a FakeResponseWriter for capturing Vektor responses
func NewFakeResponseWriter() *FakeResponseWriter {
	return &FakeResponseWriter{make(http.Header), &bytes.Buffer{}, 0}
}

// Header implements http.ResponseWriter
func (w *FakeResponseWriter) Header() http.Header {
	return w.header
}

// WriteHeader implements http.ResponseWriter
func (w *FakeResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w FakeResponseWriter) Status() int {
	return w.status
}
