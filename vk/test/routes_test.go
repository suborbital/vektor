package test_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vk/test"
)

func RoutesTest(t *testing.T) {
	server := vk.New(
		vk.UseAppName("vk tester"),
		vk.UseEnvPrefix("APP"),
		vk.UseTestMode(true),
		vk.UseInspector(func(r http.Request) {
			fmt.Println("pre-router:", r.URL.Path)
		}),
	)

	test.AddRoutes(server)
}
