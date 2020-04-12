package gapi

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

// Server represents a gust API server
type Server struct {
	*Router
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

	router := routerWithOptions(options)

	server := createGoServer(options, router)

	s := &Server{
		Router:  router,
		server:  server,
		options: options,
	}

	return s
}

// Start starts the server listening
func (s *Server) Start() error {
	// mount the root set of routes before starting
	s.mountGroup(s.Router.rootGroup())

	if s.options.AppName != "" {
		s.options.Logger.Info("starting", s.options.AppName, "...")
	}

	s.options.Logger.Info("serving on", s.server.Addr)

	if s.options.UseHTTP {
		return s.server.ListenAndServe()
	}

	return s.server.ListenAndServeTLS("", "")
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
