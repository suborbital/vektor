package main

import (
	"net/http"

	"github.com/suborbital/vektor/vk"
)

// HandleFound returns 200
func HandleFound(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	ctx.Log.Info("found!")

	return vk.R(200, "gotcha"), nil
}

// HandleNotFound returns 404
func HandleNotFound(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	return nil, vk.E(http.StatusNotFound, "Not Found")
}

// HandleMe handles Me requests
func HandleMe(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	return vk.R(200, struct{ Me string }{Me: "mario"}), nil
}

// HandleYou handles You requests
func HandleYou(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	return vk.R(201, "created, I guess"), nil
}
