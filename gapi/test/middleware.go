package main

import (
	"net/http"
	"strings"

	g "github.com/suborbital/gust/gapi"
)

func denyMiddleware(r *http.Request, ctx *g.Ctx) error {
	if strings.Contains(r.URL.Path, "hack") {
		ctx.Log.ErrorString("HACKER!!")

		return g.E(403, "begone, hacker")
	}

	return nil
}
