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

	group := g.Group("/api")
	group.GET("/me", server.With(HandleMe))

	group2 := g.Group("/v2")
	group2.GET("/you", server.With(HandleYou))

	group.AddGroup(group2)
	server.AddGroup(group)

	server.Start()
}
