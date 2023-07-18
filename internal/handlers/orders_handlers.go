package handlers

import (
	"github.com/JinFuuMugen/ya_gophermart.git/internal/database"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/logger"
	"github.com/dgrijalva/jwt-go"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "temp error", http.StatusNotFound) //TODO: change
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
		return err == nil
	}

	if !isValidOrderNumber(string(orderNumber)) {
		logger.Errorf("failed to read request body")
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		logger.Errorf("failed to get auth cookie: %v", err)
		http.Error(w, "Failed to get auth cookie", http.StatusUnauthorized)
		return
	}

	token := cookie.Value

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !parsedToken.Valid {
		logger.Errorf("invalid or expired auth token: %v", err)
		http.Error(w, "Invalid or expired auth token", http.StatusUnauthorized)
		return
	}

	username := claims["user"].(string)

	err = database.StoreOrder(string(orderNumber), username)
	if err != nil {
		logger.Errorf("failed to store order")
		http.Error(w, "Failed to save order number", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
