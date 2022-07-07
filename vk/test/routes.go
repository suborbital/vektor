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
func HandleFound(w http.ResponseWriter, _ *http.Request, ctx *vk.Ctx) error {
	ctx.Log.Info("found!")

	return vk.RespondString(ctx.Context, w, "gotcha", http.StatusOK)
}

// HandleNotFound returns 404
func HandleNotFound(_ http.ResponseWriter, _ *http.Request, _ *vk.Ctx) error {
	return vk.E(http.StatusNotFound, "Not Found")
}

// HandleMe handles Me requests
func HandleMe(w http.ResponseWriter, _ *http.Request, ctx *vk.Ctx) error {
	return vk.RespondJSON(ctx.Context, w, struct{ Me string }{Me: "mario"}, http.StatusOK)
}

// HandleYou handles You requests
func HandleYou(w http.ResponseWriter, _ *http.Request, ctx *vk.Ctx) error {
	ctx.Log.Info("calling you!")

	return vk.RespondJSON(ctx.Context, w, "created, I guess", http.StatusCreated)
}

// HandleBadMistake handles a bad mistake
func HandleBadMistake(_ http.ResponseWriter, _ *http.Request, _ *vk.Ctx) error {
	return errors.New("this is a bad idea")
}

// HandleSock hands Sock requests
func HandleSock(_ *http.Request, _ *vk.Ctx, _ *websocket.Conn) error {
	return nil
}

// HandleHTTP tests HTTP handlers
func HandleHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
