package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/logger"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/models"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var jwtSecret = []byte("gophermart")

func generateAuthToken(login string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	expirationTime := time.Now().Add(1 * time.Hour)
	claims["exp"] = expirationTime.Unix()
	claims["user"] = login

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		logger.Errorf("method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.Errorf("failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	//TODO: check if login is taken (in database)

	//http 409

	//TODO: add user to database

	token, err := generateAuthToken(user.Login)
	fmt.Println(token)
	if err != nil {
		logger.Errorf("failed to generate auth token: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600,
	}
	http.SetCookie(w, &cookie)

	response := struct {
		Message string `json:"message"`
	}{
		Message: "User registered and authenticated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Errorf("internal error while encoding response")
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Errorf("method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials models.User

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		logger.Errorf("failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// TODO: check auth in database
	// 401 (Unauthorized)

	token, err := generateAuthToken(credentials.Login)
	if err != nil {
		logger.Errorf("failed to generate auth token: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600,
	}
	http.SetCookie(w, &cookie)

	response := struct {
		Message string `json:"message"`
	}{
		Message: "User authenticated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Errorf("internal error while encoding response")
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}
