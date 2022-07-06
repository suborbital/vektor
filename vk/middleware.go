package vk

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Middleware represents a handler that runs on a request before reaching its handler
type Middleware func(HandlerFunc) HandlerFunc

// ContentTypeMiddleware allows the content-type to be set
func ContentTypeMiddleware(contentType string) Middleware {
	return func(inner HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, ctx *Ctx) error {
			ctx.RespHeaders.Set(contentTypeHeaderKey, contentType)

			return inner(w, r, ctx)
		}
	}
}

// WrapHandler takes an inner HandlerFunc, and a list of Middlewares, and returns a resolved handler that wraps the
// inner handler at the core with the passed in middlewares from first to last.
//
// For example in the following function call:
// - WrapHandler(coreHandler, panics, errors, logs, traces)
//
// The wrap would look like this:
// - incoming request -> traces -> logs -> errors -> panics -> coreHandler
func WrapHandler(handler HandlerFunc, mw ...Middleware) HandlerFunc {
	for _, m := range mw {
		if m != nil {
			handler = m(handler)
		}
	}

	return handler
}

func WrapWebsocket(handler WebSocketHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ctx *Ctx) error {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return E(http.StatusInternalServerError, err.Error())
		}

		return handler(r, ctx, conn)
	}
}

// CORSHandler enables CORS for a route
// pass "*" to allow all domains, or empty string to allow none
func CORSHandler(domain string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ctx *Ctx) error {
		enableCors(ctx, domain)

		return nil
	}
}

func enableCors(ctx *Ctx, domain string) {
	if domain != "" {
		ctx.RespHeaders.Set("Access-Control-Allow-Origin", domain)
		ctx.RespHeaders.Set("X-Requested-With", "XMLHttpRequest")
		ctx.RespHeaders.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, cache-control")
	}
}

// ErrorMiddleware returns a middleware that wraps a handler.
func ErrorMiddleware() Middleware {
	return func(inner HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, ctx *Ctx) error {
			if err := inner(w, r, ctx); err != nil {
				ctx.Log.ErrorString(fmt.Sprintf("ERROR: traceid: %s, msg: %s", ctx.RequestID(), err.Error()))

				if e, ok := err.(Error); ok {
					// we received a trusted error, which means we can pass on the status and message set on it.
					w.WriteHeader(e.Status())
					_, _ = w.Write([]byte(e.Message()))
					return nil
				}

				// we received an error from someplace else, return a generic 500
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
				return nil
			}

			return nil
		}
	}
}
