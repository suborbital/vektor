package mid

import (
	"net/http"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vlog"
)

// Logger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(log *vlog.Logger) vk.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {

		// Create the handler that will be attached in the middleware chain.
		h := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {

			log.Info("request started", "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			// Call the next handler.
			iFace, err = handler(r, ctx)

			log.Info("request completed", "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			// Return the error so it can be handled further up the chain.
			return iFace, err
		}

		return h
	}

	return m
}
