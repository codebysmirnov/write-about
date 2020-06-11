package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

type Subroute interface {
	Register(*mux.Router)
}

// ResponseJSON makes json response
func ResponseJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// RespondError makes error json format
func RespondError(w http.ResponseWriter, code int, message string) {
	ResponseJSON(w, code, map[string]string{"error": message})
}

// Rdumper show all information about request
func Rdumper(r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
}
