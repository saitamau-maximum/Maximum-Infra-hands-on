package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port        string
	DBPath      string
	SecretKey   string
	HashCost    int
	TokenExpiry time.Duration
}

func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBPath:      getEnv("DB_PATH", "database.db"),
		SecretKey:   getEnv("SECRET_KEY", "secret"),
		HashCost:    parseInt(getEnv("HASH_COST", "10")),
		TokenExpiry: paraseDuration(getEnv("TOKEN_EXPIRY", "24h")),
	}
}

func getEnv(key, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func parseInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return i
}

func paraseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	return d
}
