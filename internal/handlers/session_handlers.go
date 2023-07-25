package handlers

import (
	"encoding/json"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/database"
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
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.Errorf("failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	isTaken, err := database.CheckLoginTaken(user.Login)
	if err != nil {
		logger.Errorf("error checking login in database")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if isTaken {
		logger.Errorf("login already taken")
		http.Error(w, "Logging already taken", http.StatusConflict)
		return
	}

	err = database.StoreUser(user.Login, user.Password)
	if err != nil {
		logger.Errorf("error saving new user in database")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := generateAuthToken(user.Login)

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
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.User

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		logger.Errorf("failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var validCred bool

	validCred, err = database.UserAuth(credentials.Login, credentials.Password)

	if err != nil {
		logger.Errorf("error checking credentials")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	if !validCred {
		http.Error(w, "wrong credentials", http.StatusUnauthorized)
		return
	}

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

	w.WriteHeader(http.StatusOK)
}
