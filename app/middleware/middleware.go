package middleware

import "net/http"

type Middleware interface {
	Middleware(handler http.Handler) http.Handler
}
