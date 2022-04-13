package mid

import (
	"net/http"

	"github.com/suborbital/vektor/vk"
)

// Cors will set the response headers to allow cross-origin requests from domain.
//
// Keep in mind that you will need to handle the OPTIONS http method for the same route and pass the same CORS
// middleware to that as the browsers do a preflight check on that endpoint.
func Cors(domain string) vk.Middleware {
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		h := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			if domain != "" {
				ctx.RespHeaders.Set("Access-Control-Allow-Origin", domain)
				ctx.RespHeaders.Set("X-Requested-With", "XMLHttpRequest")
				ctx.RespHeaders.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, cache-control")
			}

			return handler(r, ctx)
		}

		return h
	}
	return m
}
