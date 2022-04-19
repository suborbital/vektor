package test

import (
	"net/http"
	"strings"

	"github.com/suborbital/vektor/vk"
)

type reqScope struct {
	ReqID  string `json:"req_id"`
	Foobar string `json:"foobar"`
}

func setScopeMiddleware() vk.Middleware {
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		f := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			scope := reqScope{
				ReqID:  ctx.RequestID(),
				Foobar: "barbaz",
			}

			ctx.UseScope(scope)

			return handler(r, ctx)
		}

		return f
	}

	return m
}

func denyMiddleware() vk.Middleware {
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		f := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			if strings.Contains(r.URL.Path, "hack") {
				ctx.Log.ErrorString("HACKER!!")
				ctx.Log.Debug("but maybe they're nice")

				return nil, vk.E(403, "begone, hacker")
			}

			return handler(r, ctx)
		}

		return f
	}

	return m
}

func headerMiddleware() vk.Middleware {
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		f := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {

			ctx.RespHeaders.Set("X-Vektor-Test", "foobar")

			return handler(r, ctx)
		}

		return f
	}

	return m
}
