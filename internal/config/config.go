package config

import (
	"os"
)

type Config struct {
	AppName   string
	Port      string
	JwtSecret string
}

func Load() Config {
	return Config{
		AppName:   os.Getenv("APP_NAME"),
		Port:      os.Getenv("PORT"),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
}
