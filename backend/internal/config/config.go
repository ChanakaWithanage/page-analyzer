package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	FetchTimeout    time.Duration
	MaxRedirects    int
	MaxBytes        int64
}

func Load() Config {
	timeoutSec := getEnvAsInt("FETCH_TIMEOUT_SECONDS", 30)
	maxRedirects := getEnvAsInt("FETCH_MAX_REDIRECTS", 5)
	maxBytes := getEnvAsInt64("FETCH_MAX_BYTES", 4<<20)

	return Config{
		Port:         getEnv("PORT", "8080"),
		FetchTimeout: time.Duration(timeoutSec) * time.Second,
		MaxRedirects: maxRedirects,
		MaxBytes:     maxBytes,
	}
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	if valStr := os.Getenv(name); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
		log.Printf("WARN: invalid int for %s, using default %d", name, defaultVal)
	}
	return defaultVal
}

func getEnvAsInt64(name string, defaultVal int64) int64 {
	if valStr := os.Getenv(name); valStr != "" {
		if val, err := strconv.ParseInt(valStr, 10, 64); err == nil {
			return val
		}
		log.Printf("WARN: invalid int64 for %s, using default %d", name, defaultVal)
	}
	return defaultVal
}
