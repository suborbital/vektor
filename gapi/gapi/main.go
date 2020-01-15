package main

import (
	g "github.com/suborbital/gust/gapi/server"
)

func main() {
	server := g.New(
		g.UseAppName("gapi tester"),
		g.UseDomain("scim.docker.cohix.ca"),
		g.UseInsecureHTTPWithEnvPort("PORT"),
	)

	server.GET("/f", server.With(HandleFound))
	server.POST("/f", server.With(HandleFound))
	server.GET("/nf", server.With(HandleNotFound))

	server.Start()
}
