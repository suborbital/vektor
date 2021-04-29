package vk

import (
	"crypto/tls"

	"github.com/suborbital/vektor/vlog"
)

// OptionsModifier takes an options struct and returns a modified Options struct
type OptionsModifier func(*Options)

// UseDomain sets the server to use a particular domain for TLS
func UseDomain(domain string) OptionsModifier {
	return func(o *Options) {
		o.Domain = domain
	}
}

// UseTLSConfig sets a TLS config that will be used for HTTPS
// This will take precedence over the Domain option in all cases
func UseTLSConfig(config *tls.Config) OptionsModifier {
	return func(o *Options) {
		o.TLSConfig = config
	}
}

// UseHTTPPort sets the HTTP port to be used:
// If domain is set, HTTP port will be used for LetsEncrypt challenge server
// If domain is NOT set, this option will put VK in insecure HTTP mode
func UseHTTPPort(port int) OptionsModifier {
	return func(o *Options) {
		o.HTTPPort = port
	}
}

// UseLogger allows a custom logger to be used
func UseLogger(logger *vlog.Logger) OptionsModifier {
	return func(o *Options) {
		o.Logger = logger
	}
}

// UseAppName allows an app name to be set (for vanity only, really....)
func UseAppName(name string) OptionsModifier {
	return func(o *Options) {
		o.AppName = name
	}
}

// UseEnvPrefix uses the provided env prefix (default VK) when looking up other options such as `VK_HTTP_PORT`
func UseEnvPrefix(prefix string) OptionsModifier {
	return func(o *Options) {
		o.EnvPrefix = prefix
	}
}
