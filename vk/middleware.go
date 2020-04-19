package vk

import (
	"net/http"

	"github.com/suborbital/gust/vlog"
)

// Middleware represents a handler that runs on a request before reaching its handler
type Middleware func(*http.Request, *Ctx) error

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

func loggerMiddleware(logger vlog.Logger) Middleware {
	return func(r *http.Request, ctx *Ctx) error {
		logger.Info(r.Method, r.URL.String())

		return nil
	}
}