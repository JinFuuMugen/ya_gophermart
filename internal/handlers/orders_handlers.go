package handlers

import (
	"github.com/JinFuuMugen/ya_gophermart.git/internal/logger"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "temp error", http.StatusNotFound)
}

func PostOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Errorf("method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		logger.Errorf("unsupported media type")
		http.Error(w, "Unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	orderNumber, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	isValidOrderNumber := func(orderNumber string) bool {

		orderNumber = strings.ReplaceAll(orderNumber, " ", "")

		_, err := strconv.ParseInt(orderNumber, 10, 64)
		if err != nil {
			return false
		}

		return true
	}

	if !isValidOrderNumber(string(orderNumber)) {
		logger.Errorf("failed to read request body")
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	//TODO: db save
	//saveOrderNumber()
	//if err != nil {
	//	http.Error(w, "failed to save order number", http.StatusInternalServerError)
	//	return
	//}

	w.WriteHeader(http.StatusAccepted)
}
