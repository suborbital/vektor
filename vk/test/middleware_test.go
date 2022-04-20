package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
	before.EXPECT().JobInMiddleware("before").Return(" [ beforeMW ran ] ").Times(4)

	before2 := &mocks.MiddlewareTester{}
	before2.EXPECT().JobInMiddleware("before").Return(" [ beforeMWII ran ] ").Times(2)

	after := &mocks.MiddlewareTester{}
	after.EXPECT().JobInMiddleware("after").Return(" [ afterMW ran ] ").Times(4)

	central := &mocks.MiddlewareTester{}
	central.EXPECT().JobInMiddleware("central").Return(" [ groupMW ran ] ").Times(1)

	// Set up handlers with the mocked middleware functions in place.
	beforeHandler := vk.WrapMiddleware([]vk.Middleware{
		mwBefore(t, before),
		mwBefore(t, before2),
		mwAfter(t, after),
	}, ctxReturnHandler)

	beforeFlipHandler := vk.WrapMiddleware([]vk.Middleware{
		mwBefore(t, before2),
		mwBefore(t, before),
		mwAfter(t, after),
	}, ctxReturnHandler)

	afterHandler := vk.WrapMiddleware([]vk.Middleware{
		mwAfter(t, after),
		mwBefore(t, before),
	}, ctxReturnHandler)

	centralHandler := vk.WrapMiddleware([]vk.Middleware{
		mwAfter(t, after),
		mwBefore(t, before),
	}, ctxReturnHandler)

	// Set up a log for the server to use.
	logger := vlog.Default(vlog.Level(vlog.LogLevelError))

	// Set up the test server.
	server := vk.New(
		vk.UseLogger(logger),
	)

	g := vk.Group("/b").Middleware(mwGroup(t, central))
	g.GET("/central", centralHandler)

	// Attach the handlers from above to the server.
	server.GET("/before", beforeHandler)
	server.GET("/beforeflip", beforeFlipHandler)
	server.GET("/after", afterHandler)
	server.AddGroup(g)

	// Create a new vtest based on the server with the handlers.
	vt := vtest.New(server)

	// Set up the requests to the two routes we've set up on the server.
	r, err := http.NewRequest(http.MethodGet, "/before", nil)
	require.NoError(t, err)

	rFlip, err := http.NewRequest(http.MethodGet, "/beforeflip", nil)
	require.NoError(t, err)

	r2, err := http.NewRequest(http.MethodGet, "/after", nil)
	require.NoError(t, err)

	rc, err := http.NewRequest(http.MethodGet, "/b/central", nil)
	require.NoError(t, err)

	t.Run("before", func(t *testing.T) {
		// Send the requests.
		resp1 := vt.Do(r, t)
		assert.Equal(t, " [ beforeMW ran ]  [ beforeMWII ran ]  [ handler ran at path /before ]  [ afterMW ran ] ", string(resp1.Body))
		assert.Equal(t, 200, resp1.Status)
	})

	t.Run("before flip", func(t *testing.T) {
		// Send the requests.
		respFlip := vt.Do(rFlip, t)
		assert.Equal(t, " [ beforeMWII ran ]  [ beforeMW ran ]  [ handler ran at path /beforeflip ]  [ afterMW ran ] ", string(respFlip.Body))
		assert.Equal(t, 200, respFlip.Status)
	})

	t.Run("after", func(t *testing.T) {
		resp2 := vt.Do(r2, t)
		assert.Equal(t, " [ beforeMW ran ]  [ handler ran at path /after ]  [ afterMW ran ] ", string(resp2.Body))
		assert.Equal(t, 200, resp2.Status)
	})

	t.Run("central", func(t *testing.T) {
		resp3 := vt.Do(rc, t)
		assert.Equal(t, " [ beforeMW ran ]  [ handler ran at path /b/central ]  [ afterMW ran ]  [ groupMW ran ] ", string(resp3.Body))
		assert.Equal(t, 200, resp3.Status)
	})

	// Check whether the middleware functions were called enough times.
	before.AssertExpectations(t)
	before2.AssertExpectations(t)
	after.AssertExpectations(t)
	central.AssertExpectations(t)
}

func ctxReturnHandler(r *http.Request, c *vk.Ctx) (interface{}, error) {
	beforeMWFromCTX := ""
	ctxValue := c.Get(testCTXKey)
	if ctxValue != nil {
		beforeMWFromCTX = ctxValue.(string)
	}

	return fmt.Sprintf("%s [ handler ran at path %s ]", beforeMWFromCTX, r.URL.Path), nil
}

func mwBefore(t *testing.T, testerMW MiddlewareTester) vk.Middleware {
	t.Helper()
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		f := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			mwResult := testerMW.JobInMiddleware("before")

			var storedValue string
			v := ctx.Get(testCTXKey)
			switch h := v.(type) {
			case string:
				storedValue = h + mwResult
			case nil:
				storedValue = mwResult
			default:
				return nil, vk.Err(http.StatusInternalServerError, "context value should have been either string or nil, got something different")
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

			responseBody := ""
			mwResult := testerMW.JobInMiddleware("after")

			switch v := iFace.(type) {
			case string:
				responseBody = v + " " + mwResult
			default:
				return nil, vk.Err(http.StatusInternalServerError, "response from handler should have been string, got something different")
			}

			return responseBody, err
		}

		return f
	}

	return m
}

func mwGroup(t *testing.T, testerMW MiddlewareTester) vk.Middleware {
	t.Helper()
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		f := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			iFace, err = handler(r, ctx)

			responseBody := ""
			mwResult := testerMW.JobInMiddleware("central")

			switch v := iFace.(type) {
			case string:
				responseBody = v + mwResult
			default:
				return nil, vk.Err(http.StatusInternalServerError, "response from handler should have been string, got something different")
			}

			return responseBody, err
		}

		return f
	}

	return m
}
