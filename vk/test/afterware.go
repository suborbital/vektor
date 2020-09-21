package main

import (
	"net/http"

	"github.com/suborbital/vektor/vk"
)

func logAfter(r *http.Request, ctx *vk.Ctx) {
	ctx.Log.Info("afterware log")
}
