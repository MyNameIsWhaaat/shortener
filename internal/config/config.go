package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	ServerTimeout time.Duration
	BaseURL       string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string
	PostgresDSN      string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	RedisAddr     string

	ShortCodeLength int
	CacheTTL        time.Duration

	RateLimitEnabled bool
	RateLimit        int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		ServerTimeout: getEnvAsDuration("SERVER_TIMEOUT", 30*time.Second),
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080"),

		PostgresHost:     getEnv("POSTGRES_HOST", "postgres"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "shortener"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "shortener"),
		PostgresDB:       getEnv("POSTGRES_DB", "shortener"),
		PostgresSSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),

		RedisHost:     getEnv("REDIS_HOST", "redis"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		ShortCodeLength: getEnvAsInt("SHORT_CODE_LENGTH", 6),
		CacheTTL:        getEnvAsDuration("CACHE_TTL", 24*time.Hour),

		RateLimitEnabled: getEnvAsBool("RATE_LIMIT_ENABLED", false),
		RateLimit:        getEnvAsInt("RATE_LIMIT", 100),
	}

	cfg.PostgresDSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
		cfg.PostgresSSLMode,
	)

	cfg.RedisAddr = fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
