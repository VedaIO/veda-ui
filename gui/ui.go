package gui

import (
	"net/http"
)

// HandlePing is a simple health check endpoint that returns a 200 OK status.
func HandlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
