package vk

// Middleware type describes a handler that wraps another handler.
type Middleware func(handlerFunc HandlerFunc) HandlerFunc

// wrapMiddleware will take a slice of middlewares and a handler, and then wrap the handler into the middlewares with
// the last middleware being the closest to the handler.
//
// mws := []Middleware{
//    error,
//    auth,
//    cors,
// }
// h := wrapMiddleware(mws, myHandler)
//
// In this instance the flow is
// Request -> error -> auth -> cors -> myHandler -> cors -> auth -> error -> Response
func wrapMiddleware(mws []Middleware, handler HandlerFunc) HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		h := mws[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
