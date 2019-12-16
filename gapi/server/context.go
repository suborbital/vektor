package server

import (
	"context"

	"github.com/julienschmidt/httprouter"
)

// Ctx serves a similar purpose to context.Context, but has some typed fields
type Ctx struct {
	context.Context
	Log    Logger
	Params httprouter.Params
}

// NewCtx creates a new Ctx
func NewCtx(log Logger, params httprouter.Params) Ctx {
	ctx := Ctx{
		Context: context.Background(),
		Log:     log,
		Params:  params,
	}

	return ctx
}
