package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/suborbital/vektor/vk"
)

const wordCountCtxKey = "dev.suborbital.wordCount"

func setupServer() *vk.Server {
	server := vk.New(vk.UseAppName("wordCount"), vk.UseHTTPPort(9090))
	api := vk.Group("/api/v1").Before(createWordCountMiddleware)
	api.POST("/wc", handlePost)

	server.AddGroup(api)

	return server
}

func createWordCountMiddleware(inner vk.HandlerFunc) vk.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ctx *vk.Ctx) (interface{}, error) {
		text, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, vk.E(http.StatusInternalServerError, err.Error())
		}

		wc := Wordcount(text)
		ctx.Set(wordCountCtxKey, wc)

		return inner(w, r, ctx)
	}
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

func handlePost(_ http.ResponseWriter, _ *http.Request, ctx *vk.Ctx) (interface{}, error) {
	wc := ctx.Get(wordCountCtxKey).(Wordcount)

	return vk.R(http.StatusOK, NewWCResponse(wc)), nil
}

func main() {
	server := setupServer()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
