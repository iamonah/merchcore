package base

import "net/http"

type HTTPHandlerWithErr func(w http.ResponseWriter, r *http.Request) error

type Middleware func(HTTPHandlerWithErr) HTTPHandlerWithErr

func Chain(h HTTPHandlerWithErr, middlewares ...Middleware) HTTPHandlerWithErr {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
