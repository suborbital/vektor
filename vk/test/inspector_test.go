package test_test

import (
	"net/http"
	"testing"

	"github.com/suborbital/vektor/vk"
	"github.com/suborbital/vektor/vlog"
	"github.com/suborbital/vektor/vtest"
)

func TestInspector(t *testing.T) {
	// suppress logging
	logger := vlog.Default(vlog.Level(vlog.LogLevelError))

	inspectorResult := ""

	server := vk.New(
		vk.UseLogger(logger),
		vk.UseInspector(func(r http.Request) {
			inspectorResult = r.URL.Path
		}),
	)

	p := "/somepath"

	server.GET(p, func(w http.ResponseWriter, r *http.Request, c *vk.Ctx) error {
		return vk.RespondWeb(c.Context, w, "", http.StatusOK)
	})

	vt := vtest.New(server)

	r, err := http.NewRequest(http.MethodGet, p, nil)

	if err != nil {
		t.Error(err)
	}

	vt.Do(r, t)

	if inspectorResult != p {
		t.Errorf("want: %s, got: %s", p, inspectorResult)
	}
}
