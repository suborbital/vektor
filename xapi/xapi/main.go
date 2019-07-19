package main

import (
	"github.com/taask/xeno/xapi"
)

func main() {
	server := x.New(
		x.UseDomain("scim.docker.cohix.ca"),
		x.UseHTTPPortFromEnv("PORT"),
	)

	server.GET("/f", x.With(HandleFound))
	server.GET("/nf", x.With(HandleNotFound))

	server.Start()
}
