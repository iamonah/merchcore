package base

import "net/http"

type Middlware func(http.HandlerFunc) http.HandlerFunc

func Chain(h http.HandlerFunc, middlwares ...Middlware) http.HandlerFunc {
	for i := len(middlwares) - 1; i >= 0; i-- {
		h = middlwares[i](h)
	}
	return h
}
