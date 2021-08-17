package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vtest"
)

// Lets us reuse the same *vk.Server instance for multiple tests and requests
var vt *vtest.VTest

func init() {
	server := vk.New(vk.UseAppName("wordcount_testing"))
	attachRoutes(server)

	vt = vtest.New(server)
}

func TestWordcount(t *testing.T) {
	body := strings.NewReader("There's a starman waiting in the sky\nHe'd like to come and meet us")

	req, err := http.NewRequest(http.MethodPost, "/api/v1/wc", body)

	if err != nil {
		t.Fatal(err)
	}

	vt.Do(req, t).
		AssertStatus(http.StatusOK).
		AssertJSON(WCResponse{
			Words:      14,
			Lines:      2,
			Characters: 66,
		})
}

func TestMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/wc", nil)

	if err != nil {
		t.Fatal(err)
	}

	vt.Do(req, t).
		AssertStatus(http.StatusMethodNotAllowed).
		AssertBodyString("Method Not Allowed\n")
}
