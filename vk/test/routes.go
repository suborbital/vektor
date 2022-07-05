package test

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/suborbital/vektor/vk"
)

// AddRoutes attaches the handlers defined below to a *vk.Server
func AddRoutes(server *vk.Server) {
	server.GET("/f", HandleFound)
	server.POST("/f", HandleFound)
	server.GET("/nf", HandleNotFound)
	server.WebSocket("/sock", HandleSock)

	v1 := vk.Group("/v1").Before(denyMiddleware, headerMiddleware)
	v1.GET("/me", HandleMe)
	v1.GET("/me/hack", HandleMe)

	v2 := vk.Group("/v2").Before(setScopeMiddleware)
	v2.GET("/you", HandleYou)
	v2.GET("/mistake", HandleBadMistake)

	api := vk.Group("/api")
	api.AddGroup(v1)
	api.AddGroup(v2)

	server.AddGroup(api)

	server.HandleHTTP(http.MethodGet, "/http", HandleHTTP)
}

// HandleFound returns 200
func HandleFound(_ http.ResponseWriter, _ *http.Request, ctx *vk.Ctx) (interface{}, error) {
	ctx.Log.Info("found!")

	return vk.R(200, "gotcha"), nil
}

// HandleNotFound returns 404
func HandleNotFound(_ http.ResponseWriter, _ *http.Request, _ *vk.Ctx) (interface{}, error) {
	return nil, vk.E(http.StatusNotFound, "Not Found")
}

// HandleMe handles Me requests
func HandleMe(_ http.ResponseWriter, _ *http.Request, _ *vk.Ctx) (interface{}, error) {
	return vk.R(200, struct{ Me string }{Me: "mario"}), nil
}

// HandleYou handles You requests
func HandleYou(_ http.ResponseWriter, _ *http.Request, ctx *vk.Ctx) (interface{}, error) {
	ctx.Log.Info("calling you!")

	return vk.R(201, "created, I guess"), nil
}

// HandleBadMistake handles a bad mistake
func HandleBadMistake(_ http.ResponseWriter, _ *http.Request, _ *vk.Ctx) (interface{}, error) {
	return nil, errors.New("this is a bad idea")
}

// HandleSock hands Sock requests
func HandleSock(_ *http.Request, _ *vk.Ctx, _ *websocket.Conn) error {
	return nil
}

// HandleHTTP tests HTTP handlers
func HandleHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
