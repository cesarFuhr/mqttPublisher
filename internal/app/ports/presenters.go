package ports

import (
	"encoding/json"
	"net/http"
)

// HTTPError Exception formatter to all http badRequests
type HTTPError struct {
	Message string `json:"message"`
}

func replyJSON(w http.ResponseWriter, code int, obj interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(obj)
}

func internalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(HTTPError{
		Message: "There was an unexpected error",
	})
}
