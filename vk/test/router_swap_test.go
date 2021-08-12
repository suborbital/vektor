package test_test

import (
	"net/http"
	"testing"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vlog"
	"github.com/suborbital/vektor/vtest"
)

func TestRouterSwap(t *testing.T) {
	// suppress logging
	logger := vlog.Default(vlog.Level(vlog.LogLevelError))

	server := vk.New(
		vk.UseLogger(logger),
		vk.UseTestMode(true),
	)

	p := "/somepath"

	server.GET(p, func(r *http.Request, c *vk.Ctx) (interface{}, error) {
		return vk.R(200, "before"), nil
	})

	vt := vtest.New(server)

	t.Run("before", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, p, nil)

		if err != nil {
			t.Error(err)
		}

		vt.Run(r, t).
			AssertBodyString("before")
	})

	newRouter := vk.NewRouter(logger)
	newRouter.GET(p, func(r *http.Request, c *vk.Ctx) (interface{}, error) {
		return vk.R(200, "after"), nil
	})

	server.SwapRouter(newRouter)

	t.Run("after", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, p, nil)

		if err != nil {
			t.Error(err)
		}

		vt.Run(r, t).
			AssertBodyString("after")
	})
}
