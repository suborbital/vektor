package test_test

import (
	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vk/test"
	"github.com/suborbital/vektor/vk/test/mocks"
	"github.com/suborbital/vektor/vlog"
	"github.com/suborbital/vektor/vtest"
	"net/http"
	"testing"
)

type testHandler struct {
	wrappedHandler http.Handler
	tester         test.RouterWrapperTester
}

func (th testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.tester.CalledIt()
	th.wrappedHandler.ServeHTTP(w, r)
}

func NewWrappedHandler(inner http.Handler, tester test.RouterWrapperTester) testHandler {
	return testHandler{
		wrappedHandler: inner,
		tester:         tester,
	}
}

func TestWrapper(t *testing.T) {
	// suppress logging
	logger := vlog.Default(vlog.Level(vlog.LogLevelError))

	rw := &mocks.RouterWrapperTester{}
	rw.On("CalledIt").Return(func() string {
		return "hello"
	}).Times(1)

	server := vk.New(
		vk.UseLogger(logger),
		vk.UseRouterWrapper(func(h http.Handler) http.Handler {
			return NewWrappedHandler(h, rw)
		}),
	)

	p := "/wrappedpath"

	server.GET(p, func(r *http.Request, c *vk.Ctx) (interface{}, error) {
		return vk.R(200, "before"), nil
	})

	vt = vtest.New(server)

	r, err := http.NewRequest(http.MethodGet, p, nil)

	if err != nil {
		t.Error(err)
	}

	vt.Do(r, t)

	rw.AssertExpectations(t)
}
