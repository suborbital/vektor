package main

import (
	"net/http"
	"strings"

	"github.com/suborbital/vektor/vk"
)

func denyMiddleware(r *http.Request, ctx *vk.Ctx) error {
	if strings.Contains(r.URL.Path, "hack") {
		ctx.Log.ErrorString("HACKER!!")

		return vk.E(403, "begone, hacker")
	}

	return nil
}

func headerMiddleware(r *http.Request, ctx *vk.Ctx) error {
	ctx.Headers.Set("X-Vektor-Test", "foobar")

	return nil
}
