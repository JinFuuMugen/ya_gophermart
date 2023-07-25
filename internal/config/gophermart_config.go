package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Addr           string `env:"RUN_ADDRESS"`
	DatabaseURI    string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{
		Addr: "localhost:8080",
	}

	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "server address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database URI")
	flag.StringVar(&cfg.AccrualAddress, "r", cfg.AccrualAddress, "accrual address")
	flag.Parse()

	if envAddr := os.Getenv("RUN_ADDRESS"); envAddr != "" {
		cfg.Addr = envAddr
	}

	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		cfg.DatabaseURI = envDatabaseDSN
	}

	if envAccrualAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualAddress != "" {
		cfg.AccrualAddress = envAccrualAddress
	}

	return cfg, nil
}
