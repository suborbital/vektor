package test_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vk/test"
	"github.com/suborbital/vektor/vlog"
	"github.com/suborbital/vektor/vtest"
)

// reuse the same TestMode server for these tests
var vt *vtest.VTest

func init() {
	// suppress logging
	logger := vlog.Default(vlog.Level(vlog.LogLevelWarn), vlog.ToFile("/dev/null"))

	server := vk.New(
		vk.UseLogger(logger),
		vk.UseAppName("vk tester"),
		vk.UseEnvPrefix("APP"),
		vk.UseTestMode(true),
	)

	test.AddRoutes(server)

	vt = vtest.New(server)
}

func TestFound(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/f", nil)

		if err != nil {
			t.Error(err)
		}

		vt.Run(r, t).
			AssertBodyString("gotcha").
			AssertStatus(200)
	})

	t.Run("POST", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodPost, "/f", nil)

		if err != nil {
			t.Error(err)
		}

		vt.Run(r, t).
			AssertBodyString("gotcha").
			AssertStatus(200)
	})
}

func TestNotFound(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/nf", nil)

	if err != nil {
		t.Error(err)
	}

	res := vk.E(http.StatusNotFound, "Not Found")
	expect, err := json.Marshal(res)

	if err != nil {
		t.Error(err)
	}

	vt.Run(r, t).
		AssertBody(expect).
		AssertStatus(404)
}

// also tests groups!
func TestMiddleware(t *testing.T) {
	t.Run("allow", func(t *testing.T) {
		me := struct{ Me string }{Me: "mario"}

		r, err := http.NewRequest(http.MethodGet, "/api/v1/me", nil)

		if err != nil {
			t.Error(err)
		}

		expect, err := json.Marshal(me)

		if err != nil {
			t.Error(err)
		}

		vt.Run(r, t).
			AssertBody(expect).
			AssertStatus(200)
	})

	t.Run("deny", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/api/v1/me/hack", nil)

		if err != nil {
			t.Error(err)
		}

		deny := vk.E(403, "begone, hacker")
		expect, err := json.Marshal(deny)

		if err != nil {
			t.Error(err)
		}

		vt.Run(r, t).
			AssertBody(expect).
			AssertStatus(403)
	})

	t.Run("header", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/api/v1/me", nil)

		if err != nil {
			t.Error(err)
		}

		vt.Run(r, t).AssertHeader("X-Vektor-Test", "foobar")
	})
}

func TestHandleHTTP(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/http", nil)

	if err != nil {
		t.Error(err)
	}

	vt.Run(r, t).
		AssertBodyString("").
		AssertStatus(204)
}

func TestBadMistake(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/api/v2/mistake", nil)

	if err != nil {
		t.Error(err)
	}

	vt.Run(r, t).
		AssertBodyString("Internal Server Error").
		AssertStatus(500)
}
