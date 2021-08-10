package test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/suborbital/vektor/vk"
)

func AddRoutes(server *vk.Server) {
	server.GET("/f", HandleFound)
	server.POST("/f", HandleFound)
	server.GET("/nf", HandleNotFound)

	v1 := vk.Group("/v1").Before(denyMiddleware, headerMiddleware)
	v1.GET("/me", HandleMe)
	v1.GET("/me/hack", HandleMe)

	v2 := vk.Group("/v2").Before(setScopeMiddleware).After(logAfter)
	v2.GET("/you", HandleYou)
	v2.GET("/mistake", HandleBadMistake)

	api := vk.Group("/api").After(getSetLogAfterware)
	api.AddGroup(v1)
	api.AddGroup(v2)

	server.AddGroup(api)

	server.HandleHTTP(http.MethodGet, "/http", HandleHTTP)
}

func main() {
	server := vk.New(
		vk.UseAppName("vk tester"),
		vk.UseEnvPrefix("APP"),
		vk.UseHTTPPort(9090),
		vk.UseInspector(func(r http.Request) {
			fmt.Println("pre-router:", r.URL.Path)
		}),
	)

	AddRoutes(server)

	// uncomment to test router swapping
	// go func() {
	// 	time.Sleep(time.Second * 5)

	// 	newRouter := vk.NewRouter(vlog.Default())
	// 	newRouter.GET("/f", HandleFound)

	// 	server.SwapRouter(newRouter)

	// 	fmt.Println("swap complete!")
	// }()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
