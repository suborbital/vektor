package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// RouteGroup represents a group of routes
type RouteGroup struct {
	prefix      string
	routes      []routeHandler
	middlewares []Middleware
}

type routeHandler struct {
	Method  string
	Path    string
	Handler httprouter.Handle
}

// Group creates a group of routes with a common prefix and middlewares
func Group(prefix string, middlewares ...Middleware) *RouteGroup {
	rg := &RouteGroup{
		prefix:      prefix,
		routes:      []routeHandler{},
		middlewares: []Middleware{},
	}

	return rg
}

// GET is a shortcut for server.Handle(http.MethodGet, path, handler)
func (g *RouteGroup) GET(path string, handler httprouter.Handle) {
	g.addRouteHandler(http.MethodGet, path, handler)
}

// HEAD is a shortcut for server.Handle(http.MethodHead, path, handler)
func (g *RouteGroup) HEAD(path string, handler httprouter.Handle) {
	g.addRouteHandler(http.MethodHead, path, handler)
}

// OPTIONS is a shortcut for server.Handle(http.MethodOptions, path, handler)
func (g *RouteGroup) OPTIONS(path string, handler httprouter.Handle) {
	g.addRouteHandler(http.MethodOptions, path, handler)
}

// POST is a shortcut for server.Handle(http.MethodPost, path, handler)
func (g *RouteGroup) POST(path string, handler httprouter.Handle) {
	g.addRouteHandler(http.MethodPost, path, handler)
}

// PUT is a shortcut for server.Handle(http.MethodPut, path, handler)
func (g *RouteGroup) PUT(path string, handler httprouter.Handle) {
	g.addRouteHandler(http.MethodPut, path, handler)
}

// PATCH is a shortcut for server.Handle(http.MethodPatch, path, handler)
func (g *RouteGroup) PATCH(path string, handler httprouter.Handle) {
	g.addRouteHandler(http.MethodPatch, path, handler)
}

// DELETE is a shortcut for server.Handle(http.MethodDelete, path, handler)
func (g *RouteGroup) DELETE(path string, handler httprouter.Handle) {
	g.addRouteHandler(http.MethodDelete, path, handler)
}

// AddGroup adds a group of routes to this group as a subgroup.
// the subgroup's prefix is added to all of the routes it contains,
// with the resulting path being "/subgroup.prefix/group.prefix/route/path/here"
func (g *RouteGroup) AddGroup(group *RouteGroup) {
	for _, rh := range group.routeHandlers() {
		fullPath := fmt.Sprintf("%s%s", ensureLeadingSlash(group.routePrefix()), ensureLeadingSlash(rh.Path))

		g.routes = append(g.routes, routeHandler{rh.Method, fullPath, rh.Handler})
	}
}

func (g *RouteGroup) routeHandlers() []routeHandler {
	return g.routes
}

func (g *RouteGroup) routePrefix() string {
	return g.prefix
}

func (g *RouteGroup) addRouteHandler(method string, path string, handler httprouter.Handle) {
	rh := routeHandler{
		Method:  method,
		Path:    path,
		Handler: handler,
	}

	g.routes = append(g.routes, rh)
}

func ensureLeadingSlash(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}
