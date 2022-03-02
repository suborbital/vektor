package test

import "net/http"

// RouterWrapperTester is a purely test interface, so we can generate a mocked implementation using mockery which then
// we can use to check whether it was called during execution in the wrapped handler.
type RouterWrapperTester interface {
	CalledIt() string
}

// Handler is a struct that wraps another handler and some extra data in the form of an implementation of the
// RouterWrapperTester interface.
type Handler struct {
	wrappedHandler http.Handler
	tester         RouterWrapperTester
}

// ServeHTTP is the method needed to implement the http.Handler interface.
func (th Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.tester.CalledIt()
	th.wrappedHandler.ServeHTTP(w, r)
}

// NewWrappedHandler is a constructor that takes in an http.Handler to wrap, and an implementation of the
// RouterWrapperTester interface. It returns a configured Handler.
func NewWrappedHandler(inner http.Handler, tester RouterWrapperTester) Handler {
	return Handler{
		wrappedHandler: inner,
		tester:         tester,
	}
}
