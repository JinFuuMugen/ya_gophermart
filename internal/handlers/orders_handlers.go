package handlers

import "net/http"

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "temp error", http.StatusNotFound)
}

func PostOrdersHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "temp error", http.StatusNotFound)
}
