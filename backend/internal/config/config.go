package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         string
	FetchTimeout time.Duration
	MaxRedirects int
	MaxBytes     int64
	EnablePprof     bool
    PprofPort       string
}

func Load() Config {
	timeoutSec := getEnvAsInt("FETCH_TIMEOUT_SECONDS", 300)
	maxRedirects := getEnvAsInt("FETCH_MAX_REDIRECTS", 5)
	maxBytes := getEnvAsInt64("FETCH_MAX_BYTES", 4<<20)

	cfg := Config{
		Port:         getEnv("PORT", "8080"),
		FetchTimeout: time.Duration(timeoutSec) * time.Second,
		MaxRedirects: maxRedirects,
		MaxBytes:     maxBytes,
		EnablePprof:  getEnv("ENABLE_PPROF", "true") == "true",
        PprofPort:    getEnv("PPROF_PORT", "6060"),
	}

	slog.Info("configuration loaded",
		"port", cfg.Port,
		"fetch_timeout", cfg.FetchTimeout,
		"max_redirects", cfg.MaxRedirects,
		"max_bytes", cfg.MaxBytes,
		"ENABLE_PPROF", cfg.EnablePprof,
		"PPROF_PORT", cfg.PprofPort,
	)

	return cfg
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
		slog.Warn("invalid int env var, using default", "key", name, "value", valStr, "default", defaultVal)
	}
	return defaultVal
}

func getEnvAsInt64(name string, defaultVal int64) int64 {
	if valStr := os.Getenv(name); valStr != "" {
		if val, err := strconv.ParseInt(valStr, 10, 64); err == nil {
			return val
		}
		slog.Warn("invalid int64 env var, using default", "key", name, "value", valStr, "default", defaultVal)
	}
	return defaultVal
}
