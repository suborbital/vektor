package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/taask/xeno/xapi"
)

func main() {
	server := xapi.New(
		xapi.UseDomain("scim.docker.cohix.ca"),
		xapi.UseHTTPPortFromEnv("PORT"),
	)

	server.GET("/f", xapi.With(func(r *http.Request, params httprouter.Params) (interface{}, error) {
		return xapi.R(200, "gotcha"), nil
	}))

	server.GET("/nf", xapi.With(func(r *http.Request, params httprouter.Params) (interface{}, error) {
		return nil, xapi.E(http.StatusNotFound, "Not Found")
	}))

	server.Start()
}
