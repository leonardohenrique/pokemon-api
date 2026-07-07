package middleware

import "net/http"

// Chain aplica múltiplos middlewares em ordem, do primeiro ao último.
// Chain(h, A, B, C) resulta em A(B(C(h))).
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
