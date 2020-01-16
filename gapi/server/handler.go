package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/suborbital/gust/glog"
)

// HandlerFunc is the gapi version of http.HandlerFinc
// instead of exposing the ResponseWriter, the function instead returns
// an object and an error, which are handled as described in `With` below
type HandlerFunc func(*http.Request, *Ctx) (interface{}, error)

// Handler handles the responses on behalf of the server
type Handler struct {
	*httprouter.Router
	middlewares []http.HandlerFunc
	getLogger   func() glog.Logger
}

//ServeHTTP serves HTTP requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check to see if the router has a handler for this path
	handler, params, _ := h.Lookup(r.Method, r.URL.Path)

	if handler != nil {
		h.getLogger().Info(r.Method, r.URL.String())
		// TODO: add middlewares here

		handler(w, r, params)
	} else {
		h.getLogger().Debug("not handled:", r.Method, r.URL.String())

		// let httprouter handle the fallthrough cases
		h.Router.ServeHTTP(w, r)
	}
}

// With returns an HTTP HandlerFunc that uses `inner` to handle the request
// inner returns a body and an error,
// the body can can be:
// - a gapi.Response object (status and body are written to w)
// - []byte (written directly to w, status 200)
// - a struct (marshalled to JSON and written to w, status 200)
//
// the error can be:
// - a gapi.Error type (status and message are written to w)
// - any other error object (status 500 and error.Error() are written to w)
//
func (h *Handler) With(inner HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var status int
		var body []byte

		ctx := NewCtx(h.getLogger(), params)

		resp, err := inner(r, ctx)
		if err != nil {
			status, body = errorOrOtherToBytes(err)
		} else {
			status, body = responseOrOtherToBytes(resp)
		}

		w.WriteHeader(status)
		w.Write(body)

		h.getLogger().Debug("handled", r.Method, r.URL.String(), fmt.Sprintf("(%d)", status))
	}
}
