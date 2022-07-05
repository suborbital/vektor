package vk

import (
	"fmt"
	"net/http"
	"strings"
)

// RouteGroup represents a group of routes
type RouteGroup struct {
	prefix     string
	httpRoutes []httpRouteHandler
	middleware []Middleware
}

type httpRouteHandler struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

// Group creates a group of routes with a common prefix and middlewares
func Group(prefix string) *RouteGroup {
	rg := &RouteGroup{
		prefix:     prefix,
		httpRoutes: []httpRouteHandler{},
		middleware: []Middleware{},
	}

	return rg
}

// GET is a shortcut for server.Handle(http.MethodGet, path, handler, middleware...)
func (g *RouteGroup) GET(path string, handler HandlerFunc, middleware ...Middleware) {
	g.Handle(http.MethodGet, path, WrapHandler(handler, middleware...))
}

// HEAD is a shortcut for server.Handle(http.MethodHead, path, handler)
func (g *RouteGroup) HEAD(path string, handler HandlerFunc, middleware ...Middleware) {
	g.Handle(http.MethodHead, path, WrapHandler(handler, middleware...))
}

// OPTIONS is a shortcut for server.Handle(http.MethodOptions, path, handler)
func (g *RouteGroup) OPTIONS(path string, handler HandlerFunc, middleware ...Middleware) {
	g.Handle(http.MethodOptions, path, WrapHandler(handler, middleware...))
}

// POST is a shortcut for server.Handle(http.MethodPost, path, handler)
func (g *RouteGroup) POST(path string, handler HandlerFunc, middleware ...Middleware) {
	g.Handle(http.MethodPost, path, WrapHandler(handler, middleware...))
}

// PUT is a shortcut for server.Handle(http.MethodPut, path, handler)
func (g *RouteGroup) PUT(path string, handler HandlerFunc, middleware ...Middleware) {
	g.Handle(http.MethodPut, path, WrapHandler(handler, middleware...))
}

// PATCH is a shortcut for server.Handle(http.MethodPatch, path, handler)
func (g *RouteGroup) PATCH(path string, handler HandlerFunc, middleware ...Middleware) {
	g.Handle(http.MethodPatch, path, WrapHandler(handler, middleware...))
}

// DELETE is a shortcut for server.Handle(http.MethodDelete, path, handler)
func (g *RouteGroup) DELETE(path string, handler HandlerFunc, middleware ...Middleware) {
	g.Handle(http.MethodDelete, path, WrapHandler(handler, middleware...))
}

// Handle adds a route to be handled
func (g *RouteGroup) Handle(method, path string, handler HandlerFunc, middleware ...Middleware) {
	g.addHttpRouteHandler(method, path, WrapHandler(handler, middleware...))
}

// WebSocket adds a websocket route to be handled.
func (g *RouteGroup) WebSocket(path string, handler WebSocketHandlerFunc) {
	g.addHttpRouteHandler(http.MethodGet, path, WrapWebSocket(handler))
}

// AddGroup adds a group of routes to this group as a subgroup.
// the subgroup's prefix is added to all of the routes it contains,
// with the resulting path being "/group.prefix/subgroup.prefix/route/path/here"
func (g *RouteGroup) AddGroup(group *RouteGroup) {
	g.httpRoutes = append(g.httpRoutes, group.httpRouteHandlers()...)
}

// Before adds middleware to the group, which are applied to every handler in the group (called before the handler)
func (g *RouteGroup) Before(middleware ...Middleware) *RouteGroup {
	return g.WithMiddlewares(middleware...)
}

// WithMiddlewares takes a list of Middlewares and will apply all of them to every handler in the group. Like in the
// WrapHandler, the first middleware is going to be the closest to each of the handlers in the group.
//
// Use this for general middlewares like logging, panic recovery, error handling, and tracing. Use the individual
// handler middlewares for endpoint specific things, like authentication.
func (g *RouteGroup) WithMiddlewares(middleware ...Middleware) *RouteGroup {
	g.middleware = append(g.middleware, middleware...)

	return g
}

// httpRouteHandlers computes the "full" path for each handler, and creates
// a HandlerFunc that chains together the group's middlewares
// before calling the inner HandlerFunc. It can be called 'recursively'
// since groups can be added to groups
func (g *RouteGroup) httpRouteHandlers() []httpRouteHandler {
	routes := make([]httpRouteHandler, len(g.httpRoutes))

	for i, r := range g.httpRoutes {
		fullPath := fmt.Sprintf("%s%s", ensureLeadingSlash(g.prefix), ensureLeadingSlash(r.Path))
		augR := httpRouteHandler{
			Method:  r.Method,
			Path:    fullPath,
			Handler: WrapHandler(r.Handler, g.middleware...),
		}

		routes[i] = augR
	}

	return routes
}

func (g *RouteGroup) addHttpRouteHandler(method string, path string, handler HandlerFunc) {
	rh := httpRouteHandler{
		Method:  method,
		Path:    path,
		Handler: handler,
	}

	g.httpRoutes = append(g.httpRoutes, rh)
}

func (g *RouteGroup) routePrefix() string {
	return g.prefix
}

func ensureLeadingSlash(path string) string {
	if path == "" {
		// handle the "root group" case
		return ""
	} else if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}
