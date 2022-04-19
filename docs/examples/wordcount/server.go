package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/suborbital/vektor/vk"
)

const wordcountCtxKey = "dev.suborbital.wordcount"

func setupServer() *vk.Server {
	server := vk.New(vk.UseAppName("wordcount"), vk.UseHTTPPort(9090))
	api := vk.Group("/api/v1").Middleware(createWordcountMiddleware())
	api.POST("/wc", handlePost)

	server.AddGroup(api)

	return server
}

func createWordcountMiddleware() vk.Middleware {
	m := func(handler vk.HandlerFunc) vk.HandlerFunc {
		f := func(r *http.Request, ctx *vk.Ctx) (iFace interface{}, err error) {
			text, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, vk.E(http.StatusInternalServerError, "reading body failed")
			}

			wc := Wordcount(text)
			ctx.Set(wordcountCtxKey, wc)

			return handler(r, ctx)
		}

		return f
	}

	return m
}

type WCResponse struct {
	Words      int `json:"words"`
	Lines      int `json:"lines"`
	Characters int `json:"chars"`
}

func NewWCResponse(wc Wordcount) *WCResponse {
	return &WCResponse{
		Words:      wc.Words(),
		Lines:      wc.Lines(),
		Characters: wc.Characters(),
	}
}

func handlePost(r *http.Request, ctx *vk.Ctx) (interface{}, error) {
	wc := ctx.Get(wordcountCtxKey).(Wordcount)

	return vk.R(http.StatusOK, NewWCResponse(wc)), nil
}

func main() {
	server := setupServer()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
