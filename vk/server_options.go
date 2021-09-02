package vk

import "net/http"

// SetGlobalOPTIONS applies the given http.Handler to all OPTIONS requests. This is useful for CORS preflight requests.
func (s *Server) SetGlobalOPTIONS(handler http.Handler) {
	if s.started.Load().(bool) {
		return
	}

	s.router.hrouter.GlobalOPTIONS = handler
}
