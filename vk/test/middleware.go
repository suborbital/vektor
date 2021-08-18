package test

import (
	"net/http"
	"strings"

	"github.com/suborbital/vektor/vk"
)

type reqScope struct {
	ReqID  string `json:"req_id"`
	Foobar string `json:"foobar"`
}

func setScopeMiddleware(r *http.Request, ctx *vk.Ctx) error {
	scope := reqScope{
		ReqID:  ctx.RequestID(),
		Foobar: "barbaz",
	}

	ctx.UseScope(scope)

	return nil
}

func denyMiddleware(r *http.Request, ctx *vk.Ctx) error {
	if strings.Contains(r.URL.Path, "hack") {
		ctx.Log.ErrorString("HACKER!!")
		ctx.Log.Debug("but maybe they're nice")

		return vk.E(403, "begone, hacker")
	}

	return nil
}

func headerMiddleware(r *http.Request, ctx *vk.Ctx) error {
	ctx.RespHeaders.Set("X-Vektor-Test", "foobar")

	return nil
}
