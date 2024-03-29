package vk

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"

	"github.com/suborbital/vektor/vlog"
)

// RouterWrapper provides a function signature to implement wrappers for routing. A good use case is to pass in an
// opentelemetry mux wrapper.
type RouterWrapper func(handler http.Handler) http.Handler

// Options are the available options for Server
type Options struct {
	AppName         string `env:"APP_NAME"`
	Domain          string `env:"DOMAIN"`
	HTTPPort        int    `env:"HTTP_PORT"`
	TLSPort         int    `env:"TLS_PORT"`
	TLSConfig       *tls.Config
	EnvPrefix       string
	QuietRoutes     []string
	Logger          *vlog.Logger
	RouterWrapper   RouterWrapper
	FallbackAddress string

	PreRouterInspector func(http.Request)
}

// defaultRouterWrapper is a default pass through option for a wrapper. This does not wrap the handler in anything.
var defaultRouterWrapper = func(innerHandler http.Handler) http.Handler {
	return innerHandler
}

func newOptsWithModifiers(mods ...OptionsModifier) *Options {
	options := &Options{}

	// set the default route wrapper to the noop pass through one.
	options.RouterWrapper = defaultRouterWrapper

	// loop through the provided options and apply the
	// modifier function to the options object
	for _, mod := range mods {
		mod(options)
	}

	envPrefix := defaultEnvPrefix
	if options.EnvPrefix != "" {
		envPrefix = options.EnvPrefix
	}

	options.finalize(envPrefix)

	return options
}

// ShouldUseTLS returns true if domain is set and/or TLS is configured
func (o *Options) ShouldUseTLS() bool {
	return o.Domain != "" || o.TLSConfig != nil
}

// HTTPPortSet returns true if the HTTP port is set
func (o *Options) HTTPPortSet() bool {
	return o.HTTPPort != 0
}

// ShouldUseHTTP returns true if insecure HTTP should be used
func (o *Options) ShouldUseHTTP() bool {
	return !o.ShouldUseTLS() && o.HTTPPortSet()
}

// finalize "locks in" the options by overriding any existing options with the version from the environment, and setting the default logger if needed
func (o *Options) finalize(prefix string) {
	// Append trailing _ if prefix is missing one
	if !strings.HasSuffix(prefix, "_") {
		prefix = prefix + "_"
	}

	if o.Logger == nil {
		o.Logger = vlog.Default(vlog.EnvPrefix(prefix))
	}

	// if no inspector was set, create an empty one
	if o.PreRouterInspector == nil {
		o.PreRouterInspector = func(_ http.Request) {}
	}

	envOpts := Options{}
	if err := envconfig.ProcessWith(context.Background(), &envOpts, envconfig.PrefixLookuper(prefix, envconfig.OsLookuper())); err != nil {
		o.Logger.Error(errors.Wrap(err, "[vk] failed to ProcessWith environment config"))
		return
	}

	o.replaceFieldsIfNeeded(&envOpts)
}

func (o *Options) replaceFieldsIfNeeded(replacement *Options) {
	if replacement.AppName != "" {
		o.AppName = replacement.AppName
	}

	if replacement.Domain != "" {
		o.Domain = replacement.Domain
	}

	if replacement.HTTPPort != 0 {
		o.HTTPPort = replacement.HTTPPort
	}

	if replacement.TLSPort != 0 {
		o.TLSPort = replacement.TLSPort
	}
}
