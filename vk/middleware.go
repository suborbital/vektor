package vk

import (
	"net/http"
)

// Middleware type describes a handler that wraps another handler.
type Middleware func(handlerFunc HandlerFunc) HandlerFunc

// BeforeWare represents a handler that runs on a request before reaching its handler
type BeforeWare func(*http.Request, *Ctx) error

// generate a HandlerFunc that passes the request through a set of BeforeWare first and Afterware after
func augmentHandler(inner HandlerFunc, middleware []BeforeWare) HandlerFunc {
	return func(r *http.Request, ctx *Ctx) (interface{}, error) {

		// run the middleware (which can error to stop progression)
		for _, m := range middleware {
			if err := m(r, ctx); err != nil {
				return nil, err
			}
		}

		resp, err := inner(r, ctx)

		return resp, err
	}
}

// wrapMiddleware will take a slice of middlewares and a handler, and then wrap the handler into the middlewares with
// the last middleware being the closest to the handler.
//
// mws := []Middleware{
//    error,
//    auth,
//    cors,
// }
// h := wrapMiddleware(mws, myHandler)
//
// In this instance the flow is
// Request -> error -> auth -> cors -> myHandler -> cors -> auth -> error -> Response
func wrapMiddleware(mws []Middleware, handler HandlerFunc) HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		h := mws[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
