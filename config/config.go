package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	GinMode     string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("GIN_MODE", "debug")

	return &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
		Port:        viper.GetString("PORT"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		GinMode:     viper.GetString("GIN_MODE"),
	}
}
