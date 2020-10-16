package vk

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/suborbital/vektor/vlog"
)

// Ctx serves a similar purpose to context.Context, but has some typed fields
type Ctx struct {
	context.Context
	Log       *vlog.Logger
	Params    httprouter.Params
	Headers   http.Header
	requestID string
	scope     interface{}
}

// NewCtx creates a new Ctx
func NewCtx(log *vlog.Logger, params httprouter.Params, headers http.Header) *Ctx {
	ctx := &Ctx{
		Context: context.Background(),
		Log:     log,
		Params:  params,
		Headers: headers,
	}

	return ctx
}

// UseScope sets an object to be the scope of the request, including setting the logger's scope
// the scope can be retrieved later with the Scope() method
func (c *Ctx) UseScope(scope interface{}) {
	c.Log = c.Log.CreateScoped(scope)

	c.scope = scope
}

// Scope retrieves the context's scope
func (c *Ctx) Scope() interface{} {
	return c.scope
}

// RequestID generates a UUID to act as a request ID and caches it on the Ctx object
func (c *Ctx) RequestID() string {
	if c.requestID == "" {
		c.requestID = uuid.New().String()
	}

	return c.requestID
}
