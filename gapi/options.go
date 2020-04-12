package gapi

import "github.com/suborbital/gust/glog"

// Options are the available options for Server
type Options struct {
	AppName  string
	Domain   string
	UseHTTP  bool
	HTTPPort string
	Logger   glog.Logger
}

func defaultOptions() Options {
	defaultOptions := Options{
		AppName:  "gust api",
		Domain:   "",
		UseHTTP:  false,
		HTTPPort: "",
		Logger:   glog.DefaultLogger(),
	}

	return defaultOptions
}
