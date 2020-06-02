package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
)

// responseJSON makes json response
func responseJSON(w http.ResponseWriter, status int, payload interface{}) {
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

// respondError makes error json format
func respondError(w http.ResponseWriter, code int, message string) {
	responseJSON(w, code, map[string]string{"error": message})
}

// rDumper show all information about request
func rDumper(r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
}
