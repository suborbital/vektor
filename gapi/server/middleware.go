package server

import "net/http"

// Middleware represents a handler that runs on a request before reaching its handler
type Middleware func(*http.Request, Ctx) error
