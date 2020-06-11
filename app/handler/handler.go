package handler

import (
	"github.com/gorilla/mux"
)

// Subroute register all model handlers
type Subroute interface {
	Register(*mux.Router)
}
