package xapi

// Options are the available options for Server
type Options struct {
	Domain   string
	UseHTTP  bool
	HTTPPort string
}

func defaultOptions() Options {
	defaultOptions := Options{
		Domain:   "",
		UseHTTP:  false,
		HTTPPort: "",
	}

	return defaultOptions
}
