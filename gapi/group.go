package gapi

import (
	"fmt"
	"net/http"
	"strings"
)

// RouteGroup represents a group of routes
type RouteGroup struct {
	prefix     string
	routes     []routeHandler
	middleware []Middleware
	afterware  []Middleware
}

type routeHandler struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

// Group creates a group of routes with a common prefix and middlewares
func Group(prefix string, middlewares ...Middleware) *RouteGroup {
	rg := &RouteGroup{
		prefix:     prefix,
		routes:     []routeHandler{},
		middleware: middlewares,
	}

	return rg
}

// GET is a shortcut for server.Handle(http.MethodGet, path, handler)
func (g *RouteGroup) GET(path string, handler HandlerFunc) {
	g.addRouteHandler(http.MethodGet, path, handler)
}

// HEAD is a shortcut for server.Handle(http.MethodHead, path, handler)
func (g *RouteGroup) HEAD(path string, handler HandlerFunc) {
	g.addRouteHandler(http.MethodHead, path, handler)
}

// OPTIONS is a shortcut for server.Handle(http.MethodOptions, path, handler)
func (g *RouteGroup) OPTIONS(path string, handler HandlerFunc) {
	g.addRouteHandler(http.MethodOptions, path, handler)
}

// POST is a shortcut for server.Handle(http.MethodPost, path, handler)
func (g *RouteGroup) POST(path string, handler HandlerFunc) {
	g.addRouteHandler(http.MethodPost, path, handler)
}

// PUT is a shortcut for server.Handle(http.MethodPut, path, handler)
func (g *RouteGroup) PUT(path string, handler HandlerFunc) {
	g.addRouteHandler(http.MethodPut, path, handler)
}

// PATCH is a shortcut for server.Handle(http.MethodPatch, path, handler)
func (g *RouteGroup) PATCH(path string, handler HandlerFunc) {
	g.addRouteHandler(http.MethodPatch, path, handler)
}

// DELETE is a shortcut for server.Handle(http.MethodDelete, path, handler)
func (g *RouteGroup) DELETE(path string, handler HandlerFunc) {
	g.addRouteHandler(http.MethodDelete, path, handler)
}

// Handle adds a route to be handled
func (g *RouteGroup) Handle(method, path string, handler HandlerFunc) {
	g.addRouteHandler(method, path, handler)
}

// AddGroup adds a group of routes to this group as a subgroup.
// the subgroup's prefix is added to all of the routes it contains,
// with the resulting path being "/subgroup.prefix/group.prefix/route/path/here"
func (g *RouteGroup) AddGroup(group *RouteGroup) {
	g.routes = append(g.routes, group.routeHandlers()...)
}

// routeHandlers computes the "full" path for each handler, and creates
// a HandlerFunc that chains together the group's middlewares
// before calling the inner HandlerFunc. It can be called 'recursively'
// since groups can be added to groups
func (g *RouteGroup) routeHandlers() []routeHandler {
	routes := make([]routeHandler, len(g.routes))

	for i, r := range g.routes {
		fullPath := fmt.Sprintf("%s%s", ensureLeadingSlash(g.prefix), ensureLeadingSlash(r.Path))
		augR := routeHandler{
			Method:  r.Method,
			Path:    fullPath,
			Handler: handlerWithMiddleware(r.Handler, g.middleware),
		}

		routes[i] = augR
	}

	return routes
}

func (g *RouteGroup) addRouteHandler(method string, path string, handler HandlerFunc) {
	rh := routeHandler{
		Method:  method,
		Path:    path,
		Handler: handler,
	}

	g.routes = append(g.routes, rh)
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
