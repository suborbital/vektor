package vk

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/suborbital/vektor/vlog"
)

// Ctx serves a similar purpose to context.Context, but has some typed fields
type Ctx struct {
	context.Context
	Log     *vlog.Logger
	Params  httprouter.Params
	Headers http.Header
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
