package vtest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// Response is a simplified analog to http.Response. Its methods can be chained together to perform
// multiple assertions on a single Response.
type Response struct {
	Body    []byte
	Status  int
	Headers http.Header
	t       *testing.T
}

// AssertStatus asserts the HTTP status code of this Response
func (r *Response) AssertStatus(status int) *Response {
	r.t.Run("status", func(t *testing.T) {
		if r.Status != status {
			t.Errorf("got status %d, want %d", r.Status, status)
		}
	})

	return r
}

// AssertHeader asserts the value of a single HTTP response header (case insenstive key)
func (r *Response) AssertHeader(key, val string) *Response {
	// normalize headers
	h := http.Header{}
	h.Add(key, val)

	return r.AssertHeaders(h)
}

// AssertHeaders asserts the values of a map of HTTP response headers
func (r *Response) AssertHeaders(headers http.Header) *Response {
	for key := range headers {
		r.t.Run(key, func(t *testing.T) {
			val := headers.Get(key)
			resv := r.Headers.Get(key)

			if resv == "" {
				r.t.Errorf("header %s: got <empty>, want '%s'", key, val)

			} else if resv != val {
				r.t.Errorf("header %s: got '%s', want '%s'", key, resv, val)
			}
		})
	}

	return r
}

// AssertBody asserts the response body is a byte-for-byte match
func (r *Response) AssertBody(body []byte) *Response {
	r.t.Run("body", func(t *testing.T) {
		if len(body) != len(r.Body) {
			t.Fatalf("body length mismatch: got %d bytes, want %d bytes", len(r.Body), len(body))
		}

		for i, v := range body {
			if v != r.Body[i] {
				t.Fatalf("byte mismatch at byte %d: got %s, want %s", i, string(r.Body[i]), string(v))
			}
		}
	})

	return r
}

const runeWindow = 25

// AssertBodyString asserts the response body is a rune-for-rune match
func (r *Response) AssertBodyString(body string) *Response {
	// pretty printing helpers for where the mismatch occurred
	trimAround := func(i int, str []rune) string {
		start := i - runeWindow
		if i < runeWindow {
			start = 0
		}

		end := i + runeWindow
		if len(str) < end {
			end = len(str)
		}

		return string(str[start:end])
	}

	context := func(i int, first, second []rune) string {
		f := trimAround(i, first)
		s := trimAround(i, second)

		offset := runeWindow
		if i < runeWindow {
			offset = i
		}

		return fmt.Sprintf("\nwant: %s\n got: %s\n      %*s", f, s, offset+1, "^")
	}

	r.t.Run("body", func(t *testing.T) {
		resRunes := []rune(string(r.Body))
		bodyRunes := []rune(body)

		if len(resRunes) != len(bodyRunes) {
			t.Fatalf("body length mismatch: got %d runes, want %d runes", len(resRunes), len(bodyRunes))
		}

		min := len(bodyRunes)
		if len(resRunes) < min {
			min = len(resRunes)
		}

		for i := 0; i < min; i++ {
			if bodyRunes[i] != resRunes[i] {
				t.Fatalf(`rune mismatch at index %d: %s`, i, context(i, bodyRunes, resRunes))
			}
		}
	})

	return r
}

// AssertJSON Marshals an interface{} to JSON and asserts that the Content-Type is set correctly and
// the response body is a byte-for-byte match
func (r *Response) AssertJSON(v interface{}) *Response {
	oldT := r.t

	r.t.Run("JSON", func(t *testing.T) {
		// We temporarily swap out r.t so that these subtests (but not later Assertions called on this
		// *Response struct) are nested together.
		r.t = t

		r.AssertHeader("Content-Type", "application/json")
		body, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}

		r.AssertBody(body)
	})

	r.t = oldT
	return r
}
