package test_test

import (
	"net/http"
	"testing"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vlog"
	"github.com/suborbital/vektor/vtest"
)

func afterwareRoute(r *http.Request, c *vk.Ctx) (interface{}, error) {
	return vk.R(200, ""), nil
}
func TestAfterware(t *testing.T) {
	// suppress logging
	logger := vlog.Default(vlog.Level(vlog.LogLevelError))

	afterwareResult := ""

	server := vk.New(
		vk.UseLogger(logger),
		vk.UseTestMode(true),
	)

	group := vk.Group("").After(func(r *http.Request, c *vk.Ctx) {
		afterwareResult = r.URL.Path
	})

	p := "/somepath"

	group.GET(p, afterwareRoute)
	server.AddGroup(group)

	vt := vtest.New(server)

	r, err := http.NewRequest(http.MethodGet, p, nil)

	if err != nil {
		t.Error(err)
	}

	vt.Run(r, t)

	if afterwareResult != p {
		t.Errorf("want: %s, got: %s", p, afterwareResult)
	}
}
