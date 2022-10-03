## Vektor has been deprecated. It is being absorbed into the [E2 Core repo](https://github.com/suborbital/e2core). No further releases will be published in this project. You can find the [docs](./docs)  for this deprecated project in this repo.

## Intro
`Vektor` enables development of modern web services in Go. Vektor is designed to simplify the development of web APIs by eliminating boilerplate, using secure defaults, providing plug-in points, and offering common pieces needed for web apps. Vektor is fairly opinionated, but aims to provide flexibility in the right places.

## Background
We see Go as the best language to build web APIs and rich backend services, and so Vektor's Go components are all focused on building those things. 

Vektor consists of components that can be used to help you build your web apps and services. Vektor components can be used alone or together.

## Components

**Vektor API (beta)**

The `vk` component is central to Vektor. It helps to quickly build production-ready API services with Go. It includes secure-by-default settings such as built-in LetsEncrypt, lots of customizability, and helpers galore. `vk` enables minimal-boilerplate servers with an intuitive wrapper around the most performant HTTP router, [httprouter](https://github.com/julienschmidt/httprouter).

**Vektor Logger (beta)**

`vlog` is a low-effort logging package that will allow for structured or text-based logging, that will easily work with popular third-party logging systems.

## Getting started

Creating a `vk` server is extremely simple:
```golang
import "github.com/suborbital/vektor/vk"

server := vk.New(
	vk.UseAppName("Vektor API Server"),
	vk.UseDomain("vektor.example.com"),
)
```
This will configure a server that sets up a LetsEncrypt certificate for `vektor.example.com` by serving content on :443 and the ACME challenge server on :80. Other options are available, see the full documentation for details.

To serve something, you'll need a handler:
```golang
type PingResponse struct {
	Ping string `json:"ping"`
}

func HandlePing(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	ctx.Log.Info("ping!")

	return PingResponse{Ping: "pong"}, nil
}
```
As you can see, handler functions don't actually concern themselves with _responding to a request_, rather just returning some data. `vk` is designed to handle the specifics of the HTTP response for you, all you need to do is return `(interface{}, error)`. Vektor handles the returned data based on a simple set of [rules](https://github.com/suborbital/vektor/blob/v0.6.0/vk/response.go). The simplest form is this: **Want to respond with JSON? Just return a struct**. To control exactly how the response behaves, check out the [Response and Error types](https://github.com/suborbital/vektor/blob/v0.6.0/vk/response.go).

Mounting handlers to the server is just as easy:
```golang
server.GET("/ping", HandlePing)
```

And finally, start your server:
```golang
if err := server.Start(); err != nil {
	log.Fatal(err)
}
```

That's just the beginning! Vektor includes powerful features like composable middleware, route groups, some handy built-in helpers, and more.

Copyright Suborbital contributors 2021.
