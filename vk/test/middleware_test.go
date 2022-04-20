package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vk/test/mocks"
	"github.com/suborbital/vektor/vlog"
	"github.com/suborbital/vektor/vtest"
)

const testCTXKey = "middlewareTesting"

func TestMiddlewares(t *testing.T) {
	// Set up mocks, expectations, and return values.
	before := &mocks.MiddlewareTester{}
	before.EXPECT().JobInMiddleware("before").Return("before returned").Times(2)

	after := &mocks.MiddlewareTester{}
	after.EXPECT().JobInMiddleware("after").Return("after returned").Times(2)

	// Set up handlers with the mocked middleware functions in place.
	beforeHandler := vk.WrapMiddleware([]vk.Middleware{
		mwBefore(t, before),
		mwAfter(t, after),
	}, func(r *http.Request, c *vk.Ctx) (interface{}, error) {
		return vk.R(200, "before"), nil
	})

	afterHandler := vk.WrapMiddleware([]vk.Middleware{
		mwAfter(t, after),
		mwBefore(t, before),
	}, func(r *http.Request, c *vk.Ctx) (interface{}, error) {
		return vk.R(200, "after"), nil
	})

	// Set up a log for the server to use.
	logger := vlog.Default(vlog.Level(vlog.LogLevelError))

	// Set up the test server.
	server := vk.New(
		vk.UseLogger(logger),
	)

	// Attach the handlers from above to the server.
	server.GET("/before", beforeHandler)
	server.GET("/after", afterHandler)

	// Create a new vtest based on the server with the handlers.
	vt := vtest.New(server)

	// Set up the requests to the two routes we've set up on the server.
	r, err := http.NewRequest(http.MethodGet, "/before", nil)
	require.NoError(t, err)

	r2, err := http.NewRequest(http.MethodGet, "/after", nil)
	require.NoError(t, err)

	// Send the requests.
	vt.Do(r, t)
	vt.Do(r2, t)

	// Check whether the middleware functions were called enough times.
	before.AssertExpectations(t)
	after.AssertExpectations(t)
}

func mwBefore(t *testing.T, testerMW MiddlewareTester) vk.Middleware {
	t.Helper()
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		f := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			_ = testerMW.JobInMiddleware("before")

			var storedValue string
			v := ctx.Get(testCTXKey)
			if v == nil {
				storedValue = "after mw going in"

			} else {
				storedValue = v.(string) + " after mw going in"
			}

			ctx.Set(testCTXKey, storedValue)

			return handler(r, ctx)
		}

		return f
	}

	return m
}

func mwAfter(t *testing.T, testerMW MiddlewareTester) vk.Middleware {
	t.Helper()
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		f := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			iFace, err = handler(r, ctx)

			_ = testerMW.JobInMiddleware("after")

			var storedValue string
			v := ctx.Get(testCTXKey)
			if v == nil {
				storedValue = "after mw going in"

			} else {
				storedValue = v.(string) + " after mw going in"
			}

			ctx.Set(testCTXKey, storedValue)

			return iFace, err
		}

		return f
	}

	return m
}
