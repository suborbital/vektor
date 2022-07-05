package vk

import (
	"net/http"
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

// CORSMiddleware enables CORS with the given domain for a route
// pass "*" to allow all domains, or empty string to allow none
func CORSMiddleware(domain string) Middleware {
	return func(inner HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, ctx *Ctx) error {
			enableCors(ctx, domain)

			return inner(w, r, ctx)
		}
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

func enableCors(ctx *Ctx, domain string) {
	if domain != "" {
		ctx.RespHeaders.Set("Access-Control-Allow-Origin", domain)
		ctx.RespHeaders.Set("X-Requested-With", "XMLHttpRequest")
		ctx.RespHeaders.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, cache-control")
	}
}
