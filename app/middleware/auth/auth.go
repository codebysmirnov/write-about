package auth

import "net/http"

// Auth interface
type Auth interface {
	Middleware(handler http.Handler) http.Handler
	Validate(token string) (bool, error)
	Generate(...Meta) (string, error)
}

// Context metadata
type Meta map[string]interface{}

const KeyUserMeta = "user"
