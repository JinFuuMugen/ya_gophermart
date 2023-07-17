package handlers

import "net/http"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "temp error", http.StatusNotFound)
}
