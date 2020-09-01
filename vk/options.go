package vk

import (
	"context"
	"fmt"
	"log"

	"github.com/sethvargo/go-envconfig"
	"github.com/suborbital/vektor/vlog"
)

// Options are the available options for Server
type Options struct {
	AppName  string `env:"APP_NAME"`
	Domain   string `env:"DOMAIN"`
	HTTPPort int    `env:"USE_HTTP_PORT"`
	Logger   *vlog.Logger
}

// ShouldUseHTTP returns true and a port string if the option is enabled
func (o Options) ShouldUseHTTP() (bool, string) {
	if o.HTTPPort != 0 {
		return true, fmt.Sprintf(":%d", o.HTTPPort)
	}

	return false, ""
}

func defaultOptions() Options {
	var o Options
	if err := envconfig.Process(context.Background(), &o); err != nil {
		log.Fatal(err)
	}

	o.Logger = vlog.Default()

	return o
}
