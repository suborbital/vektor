package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/suborbital/vektor/vk"
)

const wordcountCtxKey = "dev.suborbital.wordcount"

func attachRoutes(server *vk.Server) {
	api := vk.Group("/api/v1").Before(createWordcountMiddleware)
	api.POST("/wc", handlePost)

	server.AddGroup(api)
}

func createWordcountMiddleware(r *http.Request, ctx *vk.Ctx) error {
	text, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	wc := Wordcount(string(text))
	ctx.Set(wordcountCtxKey, wc)

	return nil
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
	server := vk.New(vk.UseAppName("wordcount"), vk.UseHTTPPort(9090))
	attachRoutes(server)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
