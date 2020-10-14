package vk

import (
	"github.com/suborbital/vektor/vlog"
)

// OptionsModifier takes an options struct and returns a modified Options struct
type OptionsModifier func(Options) Options

// UseDomain sets the server to use a particular domain for TLS
func UseDomain(domain string) OptionsModifier {
	return func(o Options) Options {
		o.Domain = domain

		return o
	}
}

// UseInsecureHTTP sets the server to serve on HTTP
func UseInsecureHTTP(port int) OptionsModifier {
	return func(o Options) Options {
		o.HTTPPort = port

		return o
	}
}

// UseLogger allows a custom logger to be used
func UseLogger(logger *vlog.Logger) OptionsModifier {
	return func(o Options) Options {
		o.Logger = logger

		return o
	}
}

// UseAppName allows an app name to be set (for vanity only, really....)
func UseAppName(name string) OptionsModifier {
	return func(o Options) Options {
		o.AppName = name

		return o
	}
}

// UseEnvPrefix uses the provided env prefix (default VK) when looking up other options such as `VK_HTTP_PORT`
func UseEnvPrefix(prefix string) OptionsModifier {
	return func(o Options) Options {
		o.EnvPrefix = prefix

		return o
	}
}
