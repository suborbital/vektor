package main

import (
	"net/http"

	g "github.com/suborbital/gust/gapi"
)

// HandleFound returns 200
func HandleFound(r *http.Request, ctx *g.Ctx) (interface{}, error) {
	ctx.Log.Info("found!")

	return g.R(200, "gotcha"), nil
}

// HandleNotFound returns 404
func HandleNotFound(r *http.Request, ctx *g.Ctx) (interface{}, error) {
	return nil, g.E(http.StatusNotFound, "Not Found")
}

// HandleMe handles Me requests
func HandleMe(r *http.Request, ctx *g.Ctx) (interface{}, error) {
	return g.R(200, "it's a me, Mario"), nil
}

// HandleYou handles You requests
func HandleYou(r *http.Request, ctx *g.Ctx) (interface{}, error) {
	return g.R(201, "created, I guess"), nil
}
