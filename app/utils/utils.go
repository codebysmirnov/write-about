package utils

import (
	"encoding/json"
	"github.com/Delisa-sama/logger"
	"net/http"
)

// Struct of error response
// TODO: add unique code
type Error struct {
	Message string
}

// ResponseJSON makes json response
func ResponseJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf("Failed to marshal response payload: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(response)
}

// RespondError makes error json format
func RespondError(w http.ResponseWriter, code int, message string) {
	ResponseJSON(w, code, Error{Message: message})
}
