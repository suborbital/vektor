package test_test

import (
	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vk/test"
	"github.com/suborbital/vektor/vk/test/mocks"
	"github.com/suborbital/vektor/vlog"
	"github.com/suborbital/vektor/vtest"
	"net/http"
)

func (vts *VektorSuite) TestWrapper() {
	// suppress logging
	logger := vlog.Default(vlog.Level(vlog.LogLevelError))

	rw := &mocks.RouterWrapperTester{}
	rw.On("CalledIt").Return(func() string {
		return "hello"
	}).Times(1)

	server := vk.New(
		vk.UseLogger(logger),
		vk.UseRouterWrapper(func(h http.Handler) http.Handler {
			return test.NewWrappedHandler(h, rw)
		}),
	)

	p := "/wrappedpath"

	server.GET(p, func(r *http.Request, c *vk.Ctx) (interface{}, error) {
		return vk.R(200, "before"), nil
	})

	vts.vt = vtest.New(server)

	r, err := http.NewRequest(http.MethodGet, p, nil)

	if err != nil {
		vts.Error(err)
	}

	vts.vt.Do(r, vts.T())

	rw.AssertExpectations(vts.T())
}
