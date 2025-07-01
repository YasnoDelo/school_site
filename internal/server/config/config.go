package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port      string `env:"PORT" default:"8080"`
	Env       string `env:"ENV" default:"production"`
	Database  string `env:"DB_URL,required"`
	StripeKey string `env:"STRIPE_KEY"`
}

func Load() *Config {
	// Загружаем .env, если он есть (dev‑среда)
	_ = godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to parse env: %v", err)
	}
	return &cfg
}
