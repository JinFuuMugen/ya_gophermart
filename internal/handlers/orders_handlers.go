package handlers

import (
	"encoding/json"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/dataaggregator"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/database"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/logger"
	"github.com/dgrijalva/jwt-go"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetOrdersHandler(addr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, _ := io.ReadAll(r.Body) //todo:remove
		logger.Warnf(string(body))    //todo:remove

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

		orders, err := dataaggregator.GetOrders(username, addr)
		if err != nil {
			logger.Errorf("error while getting orders: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ordersJSON, err := json.Marshal(orders)
		if err != nil {
			logger.Errorf("error while marshaling orders to JSON: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(ordersJSON)
	}
}

func PostOrdersHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		logger.Errorf("unsupported media type")
		http.Error(w, "Unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	orderNumber, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	isValidOrderNumber := func(orderNumber string) bool {
		orderNumber = strings.ReplaceAll(orderNumber, " ", "")
		_, err := strconv.ParseInt(orderNumber, 10, 64)
		return err == nil
	}

	if !isValidOrderNumber(string(orderNumber)) {
		logger.Errorf("invalid order number format")
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

	var code int

	code, err = database.CheckOrder(string(orderNumber), username)
	if err != nil {
		logger.Errorf("error while checking order in database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	switch code {
	case 409:
		http.Error(w, "This order was uploaded by another user", http.StatusConflict)
		return
	case 200:
		http.Error(w, "This order was already uploaded by you", http.StatusOK)
		return
	}
	err = database.StoreOrder(string(orderNumber), username)
	if err != nil {
		logger.Errorf("failed to store order: %v", err)
		http.Error(w, "Failed to save order number", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
