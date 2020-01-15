package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandlerFunc is the gapi version of http.HandlerFunc
// instead of exposing the ResponseWriter, the function instead returns
// an object and an error, which are handled as described in `With` below
type HandlerFunc func(*http.Request, httprouter.Params) (interface{}, error)

// Handler handles the responses on behalf of the server
type Handler struct {
	*httprouter.Router
	middlewares []http.HandlerFunc
}

//ServeHTTP serves HTTP requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, params, _ := h.Lookup(r.Method, r.URL.Path)

	if handler != nil {
		// TODO: add middlewares here

		handler(w, r, params)
	} else {
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
// TODO: determine if we want to use a different type for the params
func With(inner HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		resp, err := inner(r, params)
		if err != nil {
			status, body := errorOrOtherToBytes(err)
			w.WriteHeader(status)
			w.Write(body)
			return
		}

		status, body := responseOrOtherToBytes(resp)

		w.WriteHeader(status)
		w.Write(body)
	}
}
