package test_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vk/test"
	"github.com/suborbital/vektor/vlog"
	"github.com/suborbital/vektor/vtest"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestVektorSuite(t *testing.T) {
	suite.Run(t, new(VektorSuite))
}

type VektorSuite struct {
	suite.Suite
	vt *vtest.VTest
}

// Make sure that vt is reset to a base state before each test.
func (vts *VektorSuite) SetupTest() {
	// suppress logging
	logger := vlog.Default(vlog.Level(vlog.LogLevelWarn), vlog.ToFile("/dev/null"))

	server := vk.New(
		vk.UseLogger(logger),
		vk.UseAppName("vk tester"),
		vk.UseEnvPrefix("APP_"),
	)

	test.AddRoutes(server)

	vts.vt = vtest.New(server)
}

func (vts *VektorSuite) TestFound() {
	vts.Run("GET", func() {
		r, err := http.NewRequest(http.MethodGet, "/f", nil)

		if err != nil {
			vts.Error(err)
		}

		vts.vt.Do(r, vts.T()).
			AssertBodyString("gotcha").
			AssertStatus(200)
	})

	vts.Run("POST", func() {
		r, err := http.NewRequest(http.MethodPost, "/f", nil)

		if err != nil {
			vts.Error(err)
		}

		vts.vt.Do(r, vts.T()).
			AssertBodyString("gotcha").
			AssertStatus(200)
	})
}

func (vts *VektorSuite) TestNotFound() {
	r, err := http.NewRequest(http.MethodGet, "/nf", nil)

	if err != nil {
		vts.Error(err)
	}

	res := vk.E(http.StatusNotFound, "Not Found")
	expect, err := json.Marshal(res)

	if err != nil {
		vts.Error(err)
	}

	vts.vt.Do(r, vts.T()).
		AssertBody(expect).
		AssertStatus(404)
}

// also tests groups!
func (vts *VektorSuite) TestMiddleware() {
	vts.Run("allow", func() {
		me := struct{ Me string }{Me: "mario"}

		r, err := http.NewRequest(http.MethodGet, "/api/v1/me", nil)

		if err != nil {
			vts.Error(err)
		}

		expect, err := json.Marshal(me)

		if err != nil {
			vts.Error(err)
		}

		vts.vt.Do(r, vts.T()).
			AssertBody(expect).
			AssertStatus(200)
	})

	vts.Run("deny", func() {
		r, err := http.NewRequest(http.MethodGet, "/api/v1/me/hack", nil)

		if err != nil {
			vts.Error(err)
		}

		deny := vk.E(403, "begone, hacker")
		expect, err := json.Marshal(deny)

		if err != nil {
			vts.Error(err)
		}

		vts.vt.Do(r, vts.T()).
			AssertBody(expect).
			AssertStatus(403)
	})

	vts.Run("header", func() {
		r, err := http.NewRequest(http.MethodGet, "/api/v1/me", nil)

		if err != nil {
			vts.Error(err)
		}

		vts.vt.Do(r, vts.T()).AssertHeader("X-Vektor-Test", "foobar")
	})
}

func (vts *VektorSuite) TestHandleHTTP() {
	r, err := http.NewRequest(http.MethodGet, "/http", nil)

	if err != nil {
		vts.Error(err)
	}

	vts.vt.Do(r, vts.T()).
		AssertBodyString("").
		AssertStatus(204)
}

func (vts *VektorSuite) TestBadMistake() {
	r, err := http.NewRequest(http.MethodGet, "/api/v2/mistake", nil)

	if err != nil {
		vts.Error(err)
	}

	vts.vt.Do(r, vts.T()).
		AssertBodyString("Internal Server Error").
		AssertStatus(500)
}

func (vts *VektorSuite) TestSock() {
	r, err := http.NewRequest(http.MethodGet, "/sock", nil)
	r.Header.Add("Connection", "upgrade")
	r.Header.Add("Upgrade", "websocket")
	r.Header.Add("Sec-WebSocket-Version", "13")
	r.Header.Add("Sec-WebSocket-Key", "some-key")

	if err != nil {
		vts.Error(err)
	}

	vts.vt.Do(r, vts.T()).
		AssertStatus(101). // 101 Switching Protocols
		AssertHeader("Upgrade", "websocket")
}
