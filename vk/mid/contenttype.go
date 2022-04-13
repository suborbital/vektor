package mid

import (
	"net/http"

	"github.com/suborbital/vektor/vk"
)

// ContentType allows the content-type to be set.
func ContentType(contentType string) vk.Middleware {
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		h := func(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
			ctx.RespHeaders.Set(vk.ContentTypeHeaderKey, contentType)

			return handler(r, ctx)
		}

		return h
	}

	return m
}
