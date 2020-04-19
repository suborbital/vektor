package vk

import "github.com/suborbital/gust/vlog"

// Options are the available options for Server
type Options struct {
	AppName  string
	Domain   string
	UseHTTP  bool
	HTTPPort string
	Logger   vlog.Logger
}

func defaultOptions() Options {
	defaultOptions := Options{
		AppName:  "gust api",
		Domain:   "",
		UseHTTP:  false,
		HTTPPort: "",
		Logger:   vlog.DefaultLogger(),
	}

	return defaultOptions
}
