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
	mockMW := &mocks.MiddlewareTester{}
	mockMW.EXPECT().JobInMiddleware("before").Return(" [ beforeMW ran ] ").Times(4)
	mockMW.EXPECT().JobInMiddleware("beforeFlip").Return(" [ beforeFlipMW ran ] ").Times(2)
	mockMW.EXPECT().JobInMiddleware("after").Return(" [ afterMW ran ] ").Times(4)
	mockMW.EXPECT().JobInMiddleware("afterFlip").Return(" [ afterFlipMW ran ] ").Times(2)
	mockMW.EXPECT().JobInMiddleware("centralBefore").Return(" [ centralBeforeMW ran ] ").Times(2)
	mockMW.EXPECT().JobInMiddleware("centralAfter").Return(" [ centralAfterMW ran ] ").Times(2)
	mockMW.EXPECT().JobInMiddleware("centralBeforeFlip").Return(" [ centralBeforeFlipMW ran ] ").Times(2)
	mockMW.EXPECT().JobInMiddleware("centralAfterFlip").Return(" [ centralAfterFlipMW ran ] ").Times(2)

	// Set up handlers with the mocked middleware functions in place.
	beforeHandler := vk.WrapMiddleware([]vk.Middleware{
		mwBefore(t, mockMW, "before"),
		mwBefore(t, mockMW, "beforeFlip"),
		mwAfter(t, mockMW, "after"),
		mwAfter(t, mockMW, "afterFlip"),
	}, ctxReturnHandler)

	beforeFlipHandler := vk.WrapMiddleware([]vk.Middleware{
		mwBefore(t, mockMW, "beforeFlip"), // these two are flipped compared to the beforeHandler
		mwBefore(t, mockMW, "before"),
		mwAfter(t, mockMW, "afterFlip"),
		mwAfter(t, mockMW, "after"),
	}, ctxReturnHandler)

	centralHandler := vk.WrapMiddleware([]vk.Middleware{
		mwAfter(t, mockMW, "after"),
		mwBefore(t, mockMW, "before"),
	}, ctxReturnHandler)

	// Set up a log for the server to use.
	logger := vlog.Default(vlog.Level(vlog.LogLevelError))

	// Set up the test server.
	server := vk.New(
		vk.UseLogger(logger),
	)

	g := vk.Group("/b").Middleware(
		mwBefore(t, mockMW, "centralBefore"),
		mwBefore(t, mockMW, "centralBeforeFlip"),
		mwAfter(t, mockMW, "centralAfter"),
		mwAfter(t, mockMW, "centralAfterFlip"),
	)
	g.GET("/central", centralHandler)

	g2 := vk.Group("/c").Middleware(
		mwBefore(t, mockMW, "centralBeforeFlip"),
		mwBefore(t, mockMW, "centralBefore"),
		mwAfter(t, mockMW, "centralAfterFlip"),
		mwAfter(t, mockMW, "centralAfter"),
	)
	g2.GET("/central", centralHandler)

	// Attach the handlers from above to the server.
	server.GET("/before", beforeHandler)
	server.GET("/beforeflip", beforeFlipHandler)
	server.AddGroup(g)
	server.AddGroup(g2)

	// Create a new vtest based on the server with the handlers.
	vt := vtest.New(server)

	// Set up the requests to the two routes we've set up on the server.
	r, err := http.NewRequest(http.MethodGet, "/before", nil)
	require.NoError(t, err)

	rFlip, err := http.NewRequest(http.MethodGet, "/beforeflip", nil)
	require.NoError(t, err)

	rc, err := http.NewRequest(http.MethodGet, "/b/central", nil)
	require.NoError(t, err)

	rcFlip, err := http.NewRequest(http.MethodGet, "/c/central", nil)
	require.NoError(t, err)

	t.Run("before", func(t *testing.T) {
		// Send the requests.
		response := vt.Do(r, t)
		assert.Equal(t, " [ beforeMW ran ]  [ beforeFlipMW ran ]  [ handler ran at path /before ]  [ afterFlipMW ran ]  [ afterMW ran ] ", string(response.Body))
		assert.Equal(t, 200, response.Status)
	})

	t.Run("before flip", func(t *testing.T) {
		// Send the requests.
		response := vt.Do(rFlip, t)
		assert.Equal(t, " [ beforeFlipMW ran ]  [ beforeMW ran ]  [ handler ran at path /beforeflip ]  [ afterMW ran ]  [ afterFlipMW ran ] ", string(response.Body))
		assert.Equal(t, 200, response.Status)
	})

	t.Run("central", func(t *testing.T) {
		response := vt.Do(rc, t)
		assert.Equal(t, " [ centralBeforeMW ran ]  [ centralBeforeFlipMW ran ]  [ beforeMW ran ]  [ handler ran at path /b/central ]  [ afterMW ran ]  [ centralAfterFlipMW ran ]  [ centralAfterMW ran ] ", string(response.Body))
		assert.Equal(t, 200, response.Status)
	})

	t.Run("central flip", func(t *testing.T) {
		response := vt.Do(rcFlip, t)
		assert.Equal(t, " [ centralBeforeFlipMW ran ]  [ centralBeforeMW ran ]  [ beforeMW ran ]  [ handler ran at path /c/central ]  [ afterMW ran ]  [ centralAfterMW ran ]  [ centralAfterFlipMW ran ] ", string(response.Body))
		assert.Equal(t, 200, response.Status)
	})

	// Check whether the middleware functions were called enough times.
	mockMW.AssertExpectations(t)
}

func ctxReturnHandler(r *http.Request, c *vk.Ctx) (interface{}, error) {
	beforeMWFromCTX := ""
	ctxValue := c.Get(testCTXKey)
	if ctxValue != nil {
		beforeMWFromCTX = ctxValue.(string)
	}

	return fmt.Sprintf("%s [ handler ran at path %s ] ", beforeMWFromCTX, r.URL.Path), nil
}

func mwBefore(t *testing.T, testerMW MiddlewareTester, input string) vk.Middleware {
	t.Helper()
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		f := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			mwResult := testerMW.JobInMiddleware(input)

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

func mwAfter(t *testing.T, testerMW MiddlewareTester, input string) vk.Middleware {
	t.Helper()
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		f := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			iFace, err = handler(r, ctx)

			responseBody := ""
			mwResult := testerMW.JobInMiddleware(input)

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
