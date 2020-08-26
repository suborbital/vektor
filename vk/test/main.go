package main

import (
	"log"
	"net/http"

	"github.com/suborbital/vektor/vk"
)

func main() {
	server := vk.New(
		vk.UseAppName("vk tester"),
	)

	server.GET("/f", HandleFound)
	server.POST("/f", HandleFound)
	server.GET("/nf", HandleNotFound)

	v1 := vk.Group("/v1", denyMiddleware, headerMiddleware)
	v1.GET("/me", HandleMe)
	v1.GET("/me/hack", HandleMe)

	v2 := vk.Group("/v2")
	v2.GET("/you", HandleYou)
	v2.GET("/mistake", HandleBadMistake)

	api := vk.Group("/api")
	api.AddGroup(v1)
	api.AddGroup(v2)

	server.AddGroup(api)

	server.HandleHTTP(http.MethodGet, "/http", HandleHTTP)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
