package x

import "github.com/cohix/simplog"

// Options are the available options for Server
type Options struct {
	Domain   string
	UseHTTP  bool
	HTTPPort string
	Logger   Logger
}

func defaultOptions() Options {
	defaultOptions := Options{
		Domain:   "",
		UseHTTP:  false,
		HTTPPort: "",
		Logger:   simplog.New(simplog.LevelInfo),
	}

	return defaultOptions
}
