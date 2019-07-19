package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	x "github.com/taask/xeno/xapi"
)

// HandleFound returns 200
func HandleFound(r *http.Request, params httprouter.Params) (interface{}, error) {
	return x.R(200, "gotcha"), nil
}

// HandleNotFound returns 404
func HandleNotFound(r *http.Request, params httprouter.Params) (interface{}, error) {
	return nil, x.E(http.StatusNotFound, "Not Found")
}
