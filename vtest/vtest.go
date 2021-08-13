/*
Package vtest is used to create test suites for Vektor servers.

Unlike how other Go web frameworks structure their tests, the vtest package allows you to test a
Vektor server without starting an HTTP server listening on a port. Doing it this way bypasses some
of vektor's convenience features (like Let's Encrypt automation) which can get in the way in a test
environment.

Examples

This example covers basic vtest usage.

	import (
		"net/http"
		"testing"

		"github.com/suborbital/vektor/vk"
		"github.com/suborbital/vektor/vtest"
	)

	func TestExample(t *testing.T) {
		server := vk.New(
			vk.UseTestMode(true), // most vk.OptionsModifers will be ignored with TestMode enabled
		)

		handleHello := func(r *http.Request, c *vk.Ctx) (interface{}, error) {
			return vk.R(200, "Hello, Vektor!"), nil
		}

		server.GET("/hello", handleHello)

		vt := vtest.New(server)

		// vtest handles normal http.Request objects
		req, _ := http.NewRequest(http.MethodGet, "/hello", nil)

		vt.Do(req, t).
			AssertStatus(200).
			AssertBodyString("hello")
	}
*/
package vtest

import (
	"net/http"
	"testing"

	"github.com/suborbital/vektor/vk"
)

// VTest helps send normal http.Requests through the provided *vk.Server's router and creates Responses
// for testing.
type VTest struct {
	server *vk.Server
}

// New creates a VTest object and starts the test server. It is used for generating standard Go tests.
func New(server *vk.Server) *VTest {
	s := &VTest{server: server}
	s.server.Start()
	return s
}

// Do takes a normal http.Request and creates a simplified Response object. Use this to create your own tests.
func (vt *VTest) Do(req *http.Request, t *testing.T) *Response {
	wr := newFakeResponseWriter()
	vt.server.ServeHTTP(wr, req)

	return &Response{
		Body:    wr.Bytes(),
		Status:  wr.Status(),
		Headers: wr.Header(),
		t:       t,
	}
}
