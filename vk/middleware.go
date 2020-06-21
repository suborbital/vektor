package vk

import (
	"net/http"

	"github.com/suborbital/vektor/vlog"
)

// Middleware represents a handler that runs on a request before reaching its handler
type Middleware func(*http.Request, *Ctx) error

// ContentTypeMiddleware allows the content-type to be set for a handler or group
func ContentTypeMiddleware(contentType string) Middleware {
	return func(r *http.Request, ctx *Ctx) error {
		ctx.Headers.Add("Content-Type", contentType)

		return nil
	}
}

func loggerMiddleware(logger vlog.Logger) Middleware {
	return func(r *http.Request, ctx *Ctx) error {
		logger.Info(r.Method, r.URL.String())

		return nil
	}
}

// generate a HandlerFunc that passes the request through a set of Middleware first
func handlerWithMiddleware(inner HandlerFunc, middleware []Middleware) HandlerFunc {
	return func(r *http.Request, ctx *Ctx) (interface{}, error) {
		for _, m := range middleware {
			if err := m(r, ctx); err != nil {
				return nil, err
			}
		}

		return inner(r, ctx)
	}
}
