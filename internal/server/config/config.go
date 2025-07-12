package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string `env:"PORT" default:"8080"`
	Env         string `env:"ENV" default:"development"`
	DatabaseDSN string `env:"DB_URL,required"`
	SessionKey  string `env:"SESSION_KEY,required"`
}

func Load() *Config {
	_ = godotenv.Load()
	log.Println("Raw DB_URL:", os.Getenv("DB_URL"))

	return &Config{
		Port:        os.Getenv("PORT"),
		Env:         os.Getenv("ENV"),
		DatabaseDSN: os.Getenv("DB_URL"),
		SessionKey:  os.Getenv("SESSION_KEY"),
	}
}
