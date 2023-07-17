package handlers

import "net/http"

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "temp error", http.StatusNotFound)
}
