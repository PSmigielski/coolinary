// Package config reads settings from environment variables.
// Any variable X can instead be given as X_FILE pointing at a file (Docker secrets).
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServiceName    string
	GRPCPort       int
	AdminPort      int
	DatabaseURL    string
	RedisCacheURL  string
	RedisStateURL  string
	LogLevel       string
	GRPCReflection bool

	// Addresses of services this one calls (host:port).
	AIAddr       string
	NotifierAddr string
}

func Load() (Config, error) {
	c := Config{
		ServiceName:    get("SERVICE_NAME", "offers"),
		GRPCPort:       getInt("GRPC_PORT", 50051),
		AdminPort:      getInt("ADMIN_PORT", 8081),
		DatabaseURL:    get("DATABASE_URL", ""),
		RedisCacheURL:  get("REDIS_CACHE_URL", ""),
		RedisStateURL:  get("REDIS_STATE_URL", ""),
		LogLevel:       get("LOG_LEVEL", "info"),
		GRPCReflection: get("GRPC_REFLECTION", "false") == "true",
		AIAddr:         get("AI_GRPC_ADDR", "ai:50051"),
		NotifierAddr:   get("NOTIFIER_GRPC_ADDR", "notifier:50051"),
	}
	if c.DatabaseURL == "" {
		return c, fmt.Errorf("DATABASE_URL is required")
	}
	if c.RedisCacheURL == "" || c.RedisStateURL == "" {
		return c, fmt.Errorf("REDIS_CACHE_URL and REDIS_STATE_URL are required")
	}
	return c, nil
}

// get returns the env var, or the trimmed content of the file named by KEY_FILE.
func get(key, def string) string {
	if path := os.Getenv(key + "_FILE"); path != "" {
		if b, err := os.ReadFile(path); err == nil {
			return strings.TrimSpace(string(b))
		}
	}
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func getInt(key string, def int) int {
	if n, err := strconv.Atoi(get(key, "")); err == nil {
		return n
	}
	return def
}
