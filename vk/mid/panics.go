package mid

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/suborbital/vektor/vk"
)

// Panics recovers from panics and converts the panic to an error and handled in Errors.
func Panics() vk.Middleware {
	// This is the actual middleware function to be executed.
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {

		// Create the handler that will be attached in the middleware chain.
		h := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {

					// Stack trace will be provided.
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
				}
			}()

			// Call the next handler and set its return value in the err variable.
			return handler(r, ctx)
		}

		return h
	}

	return m
}
