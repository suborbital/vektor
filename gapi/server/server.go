package server

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/suborbital/gust/glog"
	"golang.org/x/crypto/acme/autocert"
)

// Server represents a gust API server
type Server struct {
	// Server embeds Handler, which embeds httprouter.Router,
	// so Server inherits all of httprouter.Router's convenience funcs,
	// but Handler controls the ServeHTTP function
	*Handler
	server  *http.Server
	options Options
}

// New creates a new gust API server
func New(opts ...OptionsModifier) *Server {
	options := defaultOptions()

	// loop through the provided options and apply the
	// modifier function to the options object
	for _, mod := range opts {
		options = mod(options)
	}

	handler := &Handler{
		Router:      httprouter.New(),
		middlewares: nil,
		getLogger: func() glog.Logger {
			return options.Logger
		},
	}

	server := createGoServer(options, handler)

	s := &Server{
		Handler: handler,
		server:  server,
		options: options,
	}

	return s
}

// Start starts the server listening
func (s *Server) Start() error {
	if s.options.AppName != "" {
		s.options.Logger.Info("starting", s.options.AppName, "...")
	}

	if s.options.UseHTTP {
		return s.server.ListenAndServe()
	}

	return s.server.ListenAndServeTLS("", "")
}

// AddGroup adds a group of handlers
func (s *Server) AddGroup(group *RouteGroup) {
	for _, r := range group.routeHandlers() {
		fullPath := fmt.Sprintf("%s%s", ensureLeadingSlash(group.routePrefix()), ensureLeadingSlash(r.Path))
		s.Handle(r.Method, fullPath, r.Handler)
	}
}

func createGoServer(options Options, handler http.Handler) *http.Server {
	if options.UseHTTP {
		return goHTTPServerWithPort(options.HTTPPort, handler)
	}

	return goTLSServerWithDomain(options.Domain, handler)
}

func goTLSServerWithDomain(domain string, handler http.Handler) *http.Server {
	m := &autocert.Manager{
		Cache:      autocert.DirCache("~/.autocert"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
	}

	go http.ListenAndServe(":80", m.HTTPHandler(nil))

	s := &http.Server{
		Addr:      ":443",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   handler,
	}

	return s
}

func goHTTPServerWithPort(port string, handler http.Handler) *http.Server {
	s := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	return s
}
