package vtest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/suborbital/vektor/vk"
)

type VKTest struct {
	server *vk.Server
}

type Response struct {
	Body    []byte
	Status  int
	Headers http.Header
}

// New creates a VKTest object and starts the test server. It is used for generating standard Go tests.
func New(server *vk.Server) *VKTest {
	s := &VKTest{server: server}
	s.server.Start()
	return s
}

// Run takes a normal http.Request and creates a simplified Response object. Use this to create your own tests.
func (vt *VKTest) Run(req *http.Request) *Response {
	wr := NewFakeResponseWriter()
	vt.server.ServeHTTP(wr, req)

	return &Response{Body: wr.Bytes(), Status: wr.Status(), Headers: wr.Header()}
}

// AssertStatus is a helper function to assert the HTTP status code of the provided request
func (vt *VKTest) AssertStatus(req *http.Request, status int) func(*testing.T) {
	return func(t *testing.T) {
		res := vt.Run(req)

		if res.Status != status {
			t.Errorf("got status %d, want %d", res.Status, status)
		}
	}
}

// AssertHeader is a helper function to assert the value of a single HTTP response header (case insenstive key)
func (vt *VKTest) AssertHeader(req *http.Request, key, val string) func(*testing.T) {
	return func(t *testing.T) {
		res := vt.Run(req)
		v := res.Headers.Get(key)

		if v == "" {
			t.Errorf("header %s: <empty>, want %s", key, val)

		} else if v != val {
			t.Errorf("header %s: %s, want %s", key, v, val)

		}
	}
}

// AssertHeaders is a helper function to assert the values of a map of HTTP response headers
func (vt *VKTest) AssertHeaders(req *http.Request, headers http.Header) func(*testing.T) {
	return func(t *testing.T) {
		res := vt.Run(req)

		for key := range headers {
			t.Run(key, func(t *testing.T) {
				val := headers.Get(key)
				resv := res.Headers.Get(key)

				if resv == "" {
					t.Errorf("header %s: got <empty>, want '%s'", key, val)

				} else if resv != val {
					t.Errorf("header %s: got '%s', want '%s'", key, resv, val)
				}
			})
		}
	}
}

// AssertBody asserts the response body is a byte-for-byte match
func (vt *VKTest) AssertBody(req *http.Request, body []byte) func(*testing.T) {
	return func(t *testing.T) {
		res := vt.Run(req)

		if len(body) != len(res.Body) {
			t.Errorf("body length mismatch: got %d, want %d", len(res.Body), len(body))
		}

		for i, v := range body {
			if v != res.Body[i] {
				t.Errorf("byte mismatch at byte %d: got %s, want %s", i, string(res.Body[i]), string(v))
			}
		}
	}
}

const RuneWindow = 25

// AssertBodyString asserts the response body is a rune-for-rune match
func (vt *VKTest) AssertBodyString(req *http.Request, body string) func(*testing.T) {
	return func(t *testing.T) {
		res := vt.Run(req)

		resRunes := []rune(string(res.Body))
		bodyRunes := []rune(body)

		min := len(bodyRunes)
		if len(resRunes) < min {
			min = len(resRunes)
		}

		// pretty printing for where the mismatch occurred
		trimAround := func(i int, str []rune) string {
			start := i - RuneWindow
			if i < RuneWindow {
				start = 0
			}

			end := i + RuneWindow
			if len(str) < end {
				end = len(str)
			}

			return string(str[start:end])
		}

		context := func(i int, first, second []rune) string {
			f := trimAround(i, first)
			s := trimAround(i, second)

			offset := RuneWindow
			if i < RuneWindow {
				offset = i
			}

			return fmt.Sprintf("\nwant: %s\n got: %s\n      %*s", f, s, offset+1, "^")
		}

		for i := 0; i < min; i++ {
			if bodyRunes[i] != resRunes[i] {
				t.Errorf(`rune mismatch at index %d: %s`, i, context(i, bodyRunes, resRunes))
				return
			}
		}
	}
}
