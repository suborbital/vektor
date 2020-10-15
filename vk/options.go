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
	AppName   string `env:"_APP_NAME"`
	Domain    string `env:"_DOMAIN"`
	HTTPPort  int    `env:"_HTTP_PORT"`
	EnvPrefix string `env:"-"`
	Logger    *vlog.Logger
}

// ShouldUseHTTP returns true and a port string if the option is enabled
func (o Options) ShouldUseHTTP() (bool, string) {
	if o.HTTPPort != 0 {
		return true, fmt.Sprintf(":%d", o.HTTPPort)
	}

	return false, ""
}

// finalize "locks in" the options by overriding any existing options with the version from the environment, and setting the default logger if needed
func (o Options) finalize(prefix string) Options {
	if err := envconfig.ProcessWith(context.Background(), &o, envconfig.PrefixLookuper(prefix, envconfig.OsLookuper())); err != nil {
		log.Fatal(err)
	}

	if o.Logger == nil {
		o.Logger = vlog.Default()
	}

	return o
}
