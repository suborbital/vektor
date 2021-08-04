package vtest_test

import (
	"net/http"
	"testing"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vlog"
	"github.com/suborbital/vektor/vtest"
)

func HandleHello(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	return vk.R(200, "hello"), nil
}

type simpleStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func HandleSimpleStruct(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	return vk.R(200, simpleStruct{"Bob", 30}), nil
}

func HandleSetHeaders(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	ctx.RespHeaders.Set("X-VK-TEST", "test")
	ctx.RespHeaders.Set("X-SUBORBITAL", "rocket launch")

	return vk.R(200, ""), nil
}

func TestVtest(t *testing.T) {
	// suppress logging
	logger := vlog.Default(vlog.Level(vlog.LogLevelError), vlog.ToFile("/dev/null"))
	server := vk.New(
		vk.UseLogger(logger),
		vk.UseTestMode(true),
	)

	server.GET("/hello", HandleHello)
	server.GET("/headers", HandleSetHeaders)
	server.GET("/simple", HandleSimpleStruct)

	vt := vtest.New(server)
	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)

	t.Run("hello", vt.AssertStatus(req, 200))
	t.Run("body", vt.AssertBodyString(req, "hello"))

	req, _ = http.NewRequest(http.MethodGet, "/headers", nil)

	headers := make(http.Header)
	headers.Add("X-VK-TEST", "test")
	headers.Add("X-SUBORBITAL", "rocket launch")

	t.Run("headers", vt.AssertHeaders(req, headers))
}
