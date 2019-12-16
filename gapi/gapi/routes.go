package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	g "github.com/suborbital/gust/gapi/server"
)

// HandleFound returns 200
func HandleFound(r *http.Request, params httprouter.Params) (interface{}, error) {
	return g.R(200, "gotcha"), nil
}

// HandleNotFound returns 404
func HandleNotFound(r *http.Request, params httprouter.Params) (interface{}, error) {
	return nil, g.E(http.StatusNotFound, "Not Found")
}
