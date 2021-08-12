// This package spawns a real HTTP server based on the test suite in the parent directory.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vk/test"
)

func main() {
	server := vk.New(
		vk.UseAppName("vk tester"),
		vk.UseEnvPrefix("APP"),
		vk.UseHTTPPort(9090),
		vk.UseInspector(func(r http.Request) {
			fmt.Println("pre-router:", r.URL.Path)
		}),
	)

	test.AddRoutes(server)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
