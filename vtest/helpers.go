package vtest

import (
	"net/http"
	"testing"

	"github.com/suborbital/vektor/vk"
)

type VKTest struct {
	server *vk.Server
}

// New creates a VKTest object and starts the test server. It is used for generating standard Go tests.
func New(server *vk.Server) *VKTest {
	s := &VKTest{server: server}
	s.server.Start()
	return s
}

// Run takes a normal http.Request and creates a simplified Response object. Use this to create your own tests.
func (vt *VKTest) Run(req *http.Request, t *testing.T) *Response {
	wr := NewFakeResponseWriter()
	vt.server.ServeHTTP(wr, req)

	return &Response{
		Body:    wr.Bytes(),
		Status:  wr.Status(),
		Headers: wr.Header(),
		t:       t,
	}
}
