package server

import (
	"crypto/tls"
	"net/http"

	"github.com/julienschmidt/httprouter"
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

	for _, mod := range opts {
		options = mod(options)
	}

	handler := &Handler{
		Router:      httprouter.New(),
		middlewares: nil,
	}

	server := createServer(options, handler)

	s := &Server{
		Handler: handler,
		server:  server,
		options: options,
	}

	return s
}

// Start starts the server listening
func (s *Server) Start() error {
	if s.options.UseHTTP {
		return s.server.ListenAndServe()
	}

	return s.server.ListenAndServeTLS("", "")
}

func createServer(options Options, handler http.Handler) *http.Server {
	if options.UseHTTP {
		return httpServerWithPort(options.HTTPPort, handler)
	}

	return tlsServerWithDomain(options.Domain, handler)
}

func tlsServerWithDomain(domain string, handler http.Handler) *http.Server {
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

func httpServerWithPort(port string, handler http.Handler) *http.Server {
	s := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	return s
}
