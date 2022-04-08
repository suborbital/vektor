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
	wsRoutes   []wsRouteHandler
	middleware []Middleware
	afterware  []Afterware
}

type httpRouteHandler struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

type wsRouteHandler struct {
	Path    string
	Handler WebSocketHandlerFunc
}

// Group creates a group of routes with a common prefix and middlewares
func Group(prefix string) *RouteGroup {
	rg := &RouteGroup{
		prefix:     prefix,
		httpRoutes: []httpRouteHandler{},
		wsRoutes:   []wsRouteHandler{},
		middleware: []Middleware{},
		afterware:  []Afterware{},
	}

	return rg
}

// GET is a shortcut for server.Handle(http.MethodGet, path, handler)
func (g *RouteGroup) GET(path string, handler HandlerFunc) {
	g.addHttpRouteHandler(http.MethodGet, path, handler)
}

// HEAD is a shortcut for server.Handle(http.MethodHead, path, handler)
func (g *RouteGroup) HEAD(path string, handler HandlerFunc) {
	g.addHttpRouteHandler(http.MethodHead, path, handler)
}

// OPTIONS is a shortcut for server.Handle(http.MethodOptions, path, handler)
func (g *RouteGroup) OPTIONS(path string, handler HandlerFunc) {
	g.addHttpRouteHandler(http.MethodOptions, path, handler)
}

// POST is a shortcut for server.Handle(http.MethodPost, path, handler)
func (g *RouteGroup) POST(path string, handler HandlerFunc) {
	g.addHttpRouteHandler(http.MethodPost, path, handler)
}

// PUT is a shortcut for server.Handle(http.MethodPut, path, handler)
func (g *RouteGroup) PUT(path string, handler HandlerFunc) {
	g.addHttpRouteHandler(http.MethodPut, path, handler)
}

// PATCH is a shortcut for server.Handle(http.MethodPatch, path, handler)
func (g *RouteGroup) PATCH(path string, handler HandlerFunc) {
	g.addHttpRouteHandler(http.MethodPatch, path, handler)
}

// DELETE is a shortcut for server.Handle(http.MethodDelete, path, handler)
func (g *RouteGroup) DELETE(path string, handler HandlerFunc) {
	g.addHttpRouteHandler(http.MethodDelete, path, handler)
}

// Handle adds a route to be handled
func (g *RouteGroup) Handle(method, path string, handler HandlerFunc) {
	g.addHttpRouteHandler(method, path, handler)
}

// WebSocket adds a websocket route to be handled
func (g *RouteGroup) WebSocket(path string, handler WebSocketHandlerFunc) {
	g.addWsRouteHandler(path, handler)
}

// AddGroup adds a group of routes to this group as a subgroup.
// the subgroup's prefix is added to all of the routes it contains,
// with the resulting path being "/group.prefix/subgroup.prefix/route/path/here"
func (g *RouteGroup) AddGroup(group *RouteGroup) {
	g.httpRoutes = append(g.httpRoutes, group.httpRouteHandlers()...)
	g.wsRoutes = append(g.wsRoutes, group.wsRouteHandlers()...)
}

// Before adds middleware to the group, which are applied to every handler in the group (called before the handler)
func (g *RouteGroup) Before(middleware ...Middleware) *RouteGroup {
	g.middleware = append(g.middleware, middleware...)

	return g
}

// After adds afterware to the group, which are applied to every handler in the group (called after the handler)
func (g *RouteGroup) After(afterware ...Afterware) *RouteGroup {
	g.afterware = append(g.afterware, afterware...)

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
			Handler: augmentHttpHandler(r.Handler, g.middleware, g.afterware),
		}

		routes[i] = augR
	}

	return routes
}

// wsRouteHandlers computes the "full" path for each handler, and creates
// a HandlerFunc that chains together the group's middlewares
// before calling the inner WebSocketHandlerFunc. It can be called 'recursively'
// since groups can be added to groups
func (g *RouteGroup) wsRouteHandlers() []wsRouteHandler {
	routes := make([]wsRouteHandler, len(g.wsRoutes))

	for i, r := range g.wsRoutes {
		fullPath := fmt.Sprintf("%s%s", ensureLeadingSlash(g.prefix), ensureLeadingSlash(r.Path))
		augR := wsRouteHandler{
			Path:    fullPath,
			Handler: augmentWsHandler(r.Handler, g.middleware, g.afterware),
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

func (g *RouteGroup) addWsRouteHandler(path string, handler WebSocketHandlerFunc) {
	rh := wsRouteHandler{
		Path:    path,
		Handler: handler,
	}

	g.wsRoutes = append(g.wsRoutes, rh)
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
