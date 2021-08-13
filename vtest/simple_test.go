package vtest_test

import (
	"net/http"
	"testing"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vlog"
	"github.com/suborbital/vektor/vtest"
)

type simpleStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func handleHello(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	return vk.R(200, "hello"), nil
}

func handleSimpleStruct(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	return vk.R(200, simpleStruct{"Bob", 30}), nil
}

func handleSetHeaders(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	ctx.RespHeaders.Set("X-VK-TEST", "test")
	ctx.RespHeaders.Set("X-SUBORBITAL", "rocket launch")

	return vk.R(200, ""), nil
}

func TestVtest(t *testing.T) {
	// suppress logging
	logger := vlog.Default(vlog.Level(vlog.LogLevelError), vlog.ToFile("/dev/null"))

	server := vk.New(
		vk.UseLogger(logger),
	)

	server.GET("/hello", handleHello)
	server.GET("/headers", handleSetHeaders)
	server.GET("/simple", handleSimpleStruct)

	vt := vtest.New(server)

	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)

	vt.Do(req, t).
		AssertStatus(200).
		AssertBodyString("hello")

	req, _ = http.NewRequest(http.MethodGet, "/headers", nil)

	t.Run("headers", func(t *testing.T) {
		headers := make(http.Header)
		headers.Add("X-VK-TEST", "test")
		headers.Add("X-SUBORBITAL", "rocket launch")

		vt.Do(req, t).AssertHeaders(headers)
	})

	req, _ = http.NewRequest(http.MethodGet, "/simple", nil)

	vt.Do(req, t).AssertStatus(200).AssertHeader("Content-Type", "application/json")
}
