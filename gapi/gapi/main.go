package main

import (
	g "github.com/suborbital/gust/gapi/server"
)

func main() {
	server := g.New(
		g.UseDomain("scim.docker.cohig.ca"),
		g.UseHTTPPortFromEnv("PORT"),
	)

	server.GET("/f", g.With(HandleFound))
	server.GET("/nf", g.With(HandleNotFound))

	server.Start()
}
