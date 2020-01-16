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

	server.GET("/f", HandleFound)
	server.POST("/f", HandleFound)
	server.GET("/nf", HandleNotFound)

	group := g.Group("/api")
	group.GET("/me", HandleMe)

	group2 := g.Group("/v2")
	group2.GET("/you", HandleYou)

	group.AddGroup(group2)
	server.AddGroup(group)

	server.Start()
}
