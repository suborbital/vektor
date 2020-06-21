# The Vektor Guide 🗺

Vektor's goal is to help you develop web services faster. Vektor handles much of the boilerplate needed to start building a Go server, so you can serve a request in less than 10 lines of code:
```golang
import "github.com/suborbital/vektor/vk"

server := vk.New(vk.UseAppName("Vektor API Server"), vk.UseDomain("vektor.example.com"))

server.GET("/ping", HandlePing)

if err := server.Start(); err != nil {
	log.Fatal(err)
}

func HandlePing(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	return "pong", nil
}
```
Those are the basics, but Vektor is capable of scaling up to serve powerful production workloads, using its full framework of API-oriented features.

# Set up `vk`

## The server object

The `vk.Server` type contains everything needed to build a web service. It includes the router, a middleware system, customizable plug-in points, and handy built-in components like LetsEncrypt support and CORS handlers.

Creating a server object is done with `vk.New()` and accepts an optional list of `OptionModifiers` which allow customization of the server:
```golang
server := vk.New(
	vk.UseAppName("Vektor API Server"),
	vk.UseDomain("vektor.example.com"),
)
```
The included `OptionsModifiers` are:

Option | Description
--- | ---
UseDomain(domain string) | Enable LetsEncrypt support with the provided domain name (will serve on :80 and :443 for challenge server and API server). LetsEncrypt is disabled by default.
UseInsecureHTTP(port int) | Choose the port on which to serve requests. Default is port 80.
UseInsecureHTTPWithEnvPort(envKey string) | Serve requests on the port set in the `envKey` environment variable.
UseLogger(logger vlog.Logger) | Set the logger object to be used, which conforms to the `vlog.Logger` interface. The logger is used internally by `vk` and is available to all handler functions via the `ctx` object. `vlog.DefaultLogger` is used by default.
UseAppName(name string) | When the application starts, `name` will be logged. Empty by default.

## Handler functions

`vk`'s handler function definition is:
```golang
func HandlePing(r *http.Request, ctx *vk.Ctx) (interface{}, error)
```
Here's a breakdown of each part:

`r *http.Request`: The request object for the request being handled.

`ctx *vk.Ctx`: A context object containing more options for interacting with the request. Ctx includes a standard Go `context.Context` which can be augmented with any value, a `vlog.Logger` object for logging within request handlers, an `httprouter.Params` object to access URL parameters (such as `/users/:uuid`), and an `http.Headers` object, which can be used to set response headers if needed.

`(interface{}, error)`: The return types of the handler allow you to respond to HTTP requests by simply returning values. If an error is returned, `vk` will interpret it as a failed request and respond with an error code, if error is `nil`, then the `interface{}` value is used to respond based on the response handling rules. **Responding to requests is handled in depth below in [Responding to requests](#responding-to-requests)**


## Mounting routes

To define routes for your `vk` server, use the HTTP method functions on the server object:
```golang
server := vk.New(
	vk.UseAppName("Vektor API Server"),
	vk.UseDomain("vektor.example.com"),
)

server.GET("/users", HandleGetUsers)
server.POST("/groups", HandleCreateGroup)
```
If you prefer to pass the HTTP method as an argument, use `server.Handle()` instead.


## Route groups

`vk` allows grouping routes by a common path prefix. For example, if you want a group of routes to begin with the `/api/` path, you can create an API route group and then mount all of your handlers to that group.
```golang
apiGroup := vk.Group("/api")
apiGroup.GET("/events", HandleGetEvents)

server.AddGroup(apiGroup)
```
Calling `AddGroup` will calculate the full paths for all routes and mount them to the server. In the example above, the handler would be mounted at `/api/events`.

Groups can even be added to groups!
```golang
v1 := vk.Group("/v1")
v1.GET("/events", HandleEventsV1)

v2 := vk.Group("/v2")
v2.GET("/events", HandleEventsV2)

apiGroup := vk.Group("/api")
apiGroup.AddGroup(v1)
apiGroup.AddGroup(v2)

server.AddGroup(api)
```
This will create a natural grouping of your routes, with the above example creating the `/api/v1/events` and `/api/v2/events` routes.


## Middleware

Groups become even more powerful when combined with Middleware. Middleware are pseudo request handlers that run in sequence before the mounted `vk.HandlerFunc` is run. Middleware functions can modify a request and its context, or they can return an error, which causes the request handling to be terminated immediately. Two examples:
```golang
func headerMiddleware(r *http.Request, ctx *vk.Ctx) error {
	ctx.Headers.Set("X-Vektor-Test", "foobar")

	return nil
}

func denyMiddleware(r *http.Request, ctx *vk.Ctx) error {
	if strings.Contains(r.URL.Path, "hack") {
		ctx.Log.ErrorString("HACKER!!")

		return vk.E(403, "begone, hacker")
	}

	return nil
}
```
As you can see, middleware have a similar function signature to `vk.HandlerFunc`, but only return an error. The first example modifies the request context to add a response header. The second example detects a hacker and returns an error, which is handled exactly like any other error response (see below). Returning an error from a Middleware prevents the request from ever reaching the registered handler.

Middleware are applied to route groups:
```golang
v1 := vk.Group("/v1", vk.ContentTypeMiddleware("application/json"), denyMiddleware, headerMiddleware)
v1.GET("/events", HandleEventsV1)
```
This example shows a group created with three middleware. The first adds the `Content-Type` response header (and is included with `vk`), the second and third are the examples from above. When the group is mounted to the server, the chain of middleware are put in place, and are run before the registered handler. When groups are nested, the middleware from the parent group are run before the middleware of any child groups. In the example of nested groups above, any middleware set on the `apiGroup` groups would run before any middleware set on the `v1` or `v2` groups.

Middleware in `vk` is designed to be easily composable, creating chains of behaviour easily grouped to sets of routes. Middleware can also help increase security of applications, allowing authentication, request throttling, active defence, etc, to run before the registered handler and keeping sensitive code from even being reached in the case of an unauthorized request.


# Responding to requests

## Response types

`vk` includes two types, `Response` and `Error` (with helper functions `vk.Respond` and `vk.Error`) that can be used to gain extra control over the response code and contents that you want to return:

```golang
type createdResponse struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

func HandleCreate(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	// Do some work

	resp := createdResponse {
		Name: "Wendy",
		UUID: "jfioqerjhp98qergnjw98h23"
	}

	// Return 201 (Created) and JSON
	return vk.Respond(http.StatusCreated, resp), nil
}

func HandleDelete(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	// Oops, something went wrong

	return nil, vk.Error(http.StatusConflict, "the user is already deleted") // responds with HTTP status 409 and body {"status": 409, "message": "the user is already deleted"}
}
```
`vk.Respond` and `vk.Error` can be used with their shortcuts `vk.R` and `vk.E` if you like your code to be terse.

## Response handling rules

`vk` processes the `(interface{}, error)` returned by handler functions in a spcific way to ensure you always know how it will behave while still being able to use simple types in your code.

### Successful responses (i.e. the `interface{}` returned by handler functions):

1. If the type is `vk.Response`, set the HTTP response code provided and process `Response.body` as follows. (If the type is NOT `vk.Response`, the status code is set to `200 OK`)
1. If the type is string, write the string (as UTF-8 bytes) to the response body.
1. If the type is bytes, write them directly to the response body.
1. If the type is a struct, attempt to marshal to JSON and write JSON bytes to the response body.

Examples:

Handler returns... | Status Code | Response body
--- | --- | ---
`return "Hello, World", nil` | 200 OK | "Hello World" (as UTF-8 bytes)
`return jsonBytesFromJSONMarshal, nil` | 200 OK | [JSON bytes as generated by json.Marshal]
`return someStructInstance, nil` | 200 OK | [JSON respresentation of struct automatically marshalled by `vk`]
`return vk.R(http.StatusCreated, "created"), nil` | 201 Created | "created" (as UTF-8 bytes)
`return vk.R(http.StatusCreated, someStructInstance), nil` | 201 Created | [JSON respresentation of struct automatically marshalled by `vk`]

### Failure responses (i.e. the `error` returned by handler functions):

1. If the type is `vk.Error`, set the HTTP response code provided and respond with JSON as follows: `{"status": $code, "message": $message}
2. If the type is NOT `vk.Error`, set the HTTP status code to 500 and write `err.Error()` as UTF-8 bytes to the response body

Examples:

Handler returns... | Status Code | Response body
--- | --- | ---
`return nil, errors.New("failed to add user")` | 500 Internal Server Error | "failed to add user" (as UTF-8 bytes)
`retuen nil, vk.E(http.StatusForbidden, "not permitted to do this thing")` | 403 Forbidden | `{"status": 403, "message": "not permitted to do this thing"}`

## What's to come?

`Vektor` is under active development. It intertwines closely with [Hive](https://github.com/suborbital/hive) to achieve Suborbital's goal of creating a framework for scalable web services. Hive and Vektor together can handle very large scale systems, and will be further integrated together to enable FaaS, WASM-based web service logic, and vastly improved developer experience and productivity.