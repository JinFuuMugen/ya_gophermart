package main

import (
	"github.com/JinFuuMugen/ya_gophermart.git/internal/config"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/database"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/handlers"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/logger"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("cannot create config: %s", err)
	}

	if err := logger.Init(); err != nil {
		log.Fatalf("cannot create logger: %s", err)
	}

	if cfg.DatabaseURI != "" {
		err := database.InitDatabase(cfg.DatabaseURI)
		if err != nil {
			logger.Fatalf("cannot create database connection: %s", err)
		}
	}

	rout := chi.NewRouter()

	rout.Route("/api/user", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return logger.HandlerLogger(next)
		})

		r.Post("/register", handlers.RegisterHandler) //TODO: make this handler
		r.Post("/login", handlers.LoginHandler)       //TODO: make this handler

		r.With(handlers.AuthMiddleware).Group(func(r chi.Router) { //TODO: make this middleware
			r.Post("/orders", handlers.PostOrdersHandler)             //TODO: make this handler
			r.Get("/orders", handlers.GetOrdersHandler)               //TODO: make this handler
			r.Get("/balance", handlers.GetBalanceHandler)             //TODO: make this handler
			r.Post("/balance/withdraw", handlers.PostWithdrawHandler) //TODO: make this handler
			r.Get("/withdrawals", handlers.GetWithdrawalsHandler)     //TODO: make this handler
		})
	})

	err = http.ListenAndServe(cfg.Addr, rout)
	if err != nil {
		logger.Fatalf("cannot start server: %s", err)
	}
}
