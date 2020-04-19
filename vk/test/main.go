package main

import (
	g "github.com/suborbital/vektor/vk"
)

func main() {
	server := g.New(
		g.UseAppName("vk tester"),
		g.UseDomain("scim.docker.cohix.ca"),
		g.UseInsecureHTTPWithEnvPort("PORT"),
	)

	server.GET("/f", HandleFound)
	server.POST("/f", HandleFound)
	server.GET("/nf", HandleNotFound)

	group := g.Group("/api")

	group3 := g.Group("/v1", denyMiddleware)
	group3.GET("/*name", HandleMe)
	group.AddGroup(group3)

	group2 := g.Group("/v2")
	group2.GET("/you", HandleYou)
	group.AddGroup(group2)

	server.AddGroup(group)

	server.Start()
}
