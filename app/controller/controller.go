package controller

import (
	"github.com/gorilla/mux"
)

// Controller register all controller handlers
type Controller interface {
	Register(*mux.Router)
}
