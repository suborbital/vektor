package test

import (
	"net/http"

	"github.com/suborbital/vektor/vk"
)

func logAfter(r *http.Request, ctx *vk.Ctx) {
	ctx.Log.Info("afterware log")
}

func getSetLogAfterware(r *http.Request, ctx *vk.Ctx) {
	ctx.Set("logMe", "loggexander contexton")

	val := ctx.Get("logMe").(string)

	ctx.Log.Debug(val)
	ctx.Log.Warn(val)
}
