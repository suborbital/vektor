package server

import (
	"context"

	"github.com/julienschmidt/httprouter"

	"github.com/suborbital/gust/glog"
)

// Ctx serves a similar purpose to context.Context, but has some typed fields
type Ctx struct {
	context.Context
	Log    glog.Logger
	Params httprouter.Params
}

// NewCtx creates a new Ctx
func NewCtx(log glog.Logger, params httprouter.Params) Ctx {
	ctx := Ctx{
		Context: context.Background(),
		Log:     log,
		Params:  params,
	}

	return ctx
}
