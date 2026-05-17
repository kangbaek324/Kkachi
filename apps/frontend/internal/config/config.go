package config

import "os"

type Config struct {
	APIBaseURL string
}

func Load() Config {
	return Config{
		APIBaseURL: getEnv("KKACHI_API_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
