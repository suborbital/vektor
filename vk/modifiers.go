package vk

import (
	"fmt"
	"os"

	"github.com/suborbital/gust/vlog"
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
		o.Domain = ""
		o.UseHTTP = true
		o.HTTPPort = fmt.Sprintf(":%d", port)

		return o
	}
}

// UseInsecureHTTPWithEnvPort sets the server to run in
// insecure HTTP mode serving on the port
// indicated by the `key` env var,
// **only if it is set**
func UseInsecureHTTPWithEnvPort(envKey string) OptionsModifier {
	return func(o Options) Options {
		port, usePort := os.LookupEnv(envKey)
		if usePort {
			o.Domain = ""
			o.UseHTTP = true
			o.HTTPPort = fmt.Sprintf(":%s", port)
		}

		return o
	}
}

// UseLogger allows a custom logger to be used
func UseLogger(logger vlog.Logger) OptionsModifier {
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
