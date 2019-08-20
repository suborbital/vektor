package x

import (
	"fmt"
	"os"
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

// UseHTTPPortFromEnv sets the server to run in
// insecure HTTP mode serving on the port
// indicated by the `key` env var,
// **only if it is set**
func UseHTTPPortFromEnv(key string) OptionsModifier {
	return func(o Options) Options {
		port, usePort := os.LookupEnv(key)
		if usePort {
			o.Domain = ""
			o.UseHTTP = true
			o.HTTPPort = fmt.Sprintf(":%s", port)
		}

		return o
	}
}

// UseLogger provides a new logger to be used
func UseLogger(l Logger) OptionsModifier {
	return func(o Options) Options {
		o.Logger = l

		return o
	}
}
