package vtest

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
	"strconv"
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

// Write implements http.ResponseWriter
func (w *fakeResponseWriter) Write(content []byte) (int, error) {
	lines := bytes.Split(content, []byte("\r\n"))

	parseHeaders := false

	// Very crude implementation of an HTTP wire format parser
	for _, line := range lines {
		// If the line is the protocol line, set the status and start looking
		// for headers
		if bytes.HasPrefix(line, []byte("HTTP")) {
			_, after, _ := bytes.Cut(line, []byte(" "))
			status, err := strconv.Atoi(string(after[:3]))
			if err != nil {
				panic(err)
			}
			w.status = status
			parseHeaders = true
			continue
		}

		// Per the spec, if we hit a blank line the next line is the body
		if len(line) == 0 {
			parseHeaders = false
			continue
		}

		// Look for headers or just write the body
		if parseHeaders {
			key, value, _ := bytes.Cut(line, []byte(": "))
			w.header[string(key)] = []string{string(value)}
		} else {
			if _, err := w.Buffer.Write(line); err != nil {
				return 0, err
			}
			break
		}
	}

	return len(content), nil
}

func (w fakeResponseWriter) Status() int {
	return w.status
}

// Hijack implements the Hijacker iterface
func (w *fakeResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rw := &bufio.ReadWriter{
		Reader: bufio.NewReader(w.Buffer),
		Writer: bufio.NewWriter(w.Buffer),
	}

	return &fakeConn{ResponseWriter: w}, rw, nil
}
