package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Save original env vars
	originalPort := os.Getenv("SERVER_PORT")
	originalTimeout := os.Getenv("SERVER_TIMEOUT")
	originalBaseURL := os.Getenv("BASE_URL")

	defer func() {
		// Restore env vars
		if originalPort != "" {
			os.Setenv("SERVER_PORT", originalPort)
		} else {
			os.Unsetenv("SERVER_PORT")
		}
		if originalTimeout != "" {
			os.Setenv("SERVER_TIMEOUT", originalTimeout)
		} else {
			os.Unsetenv("SERVER_TIMEOUT")
		}
		if originalBaseURL != "" {
			os.Setenv("BASE_URL", originalBaseURL)
		} else {
			os.Unsetenv("BASE_URL")
		}
	}()

	t.Run("default values", func(t *testing.T) {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_TIMEOUT")
		os.Unsetenv("BASE_URL")

		cfg := Load()

		if cfg.ServerPort != "8080" {
			t.Errorf("expected default ServerPort 8080, got %s", cfg.ServerPort)
		}

		if cfg.ServerTimeout != 30*time.Second {
			t.Errorf("expected default ServerTimeout 30s, got %v", cfg.ServerTimeout)
		}

		if cfg.BaseURL != "http://localhost:8080" {
			t.Errorf("expected default BaseURL, got %s", cfg.BaseURL)
		}
	})

	t.Run("custom values", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "9000")
		os.Setenv("SERVER_TIMEOUT", "60s")
		os.Setenv("BASE_URL", "https://short.example.com")

		cfg := Load()

		if cfg.ServerPort != "9000" {
			t.Errorf("expected ServerPort 9000, got %s", cfg.ServerPort)
		}

		if cfg.ServerTimeout != 60*time.Second {
			t.Errorf("expected ServerTimeout 60s, got %v", cfg.ServerTimeout)
		}

		if cfg.BaseURL != "https://short.example.com" {
			t.Errorf("expected custom BaseURL, got %s", cfg.BaseURL)
		}
	})

	t.Run("PostgreSQL DSN", func(t *testing.T) {
		os.Setenv("POSTGRES_USER", "testuser")
		os.Setenv("POSTGRES_PASSWORD", "testpass")
		os.Setenv("POSTGRES_HOST", "localhost")
		os.Setenv("POSTGRES_PORT", "5432")
		os.Setenv("POSTGRES_DB", "testdb")

		cfg := Load()

		expectedDSN := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
		if cfg.PostgresDSN != expectedDSN {
			t.Errorf("expected DSN %s, got %s", expectedDSN, cfg.PostgresDSN)
		}
	})

	t.Run("Redis config", func(t *testing.T) {
		os.Setenv("REDIS_HOST", "redis")
		os.Setenv("REDIS_PORT", "6379")
		os.Setenv("REDIS_PASSWORD", "mypassword")
		os.Setenv("REDIS_DB", "0")

		cfg := Load()

		if cfg.RedisAddr != "redis:6379" {
			t.Errorf("expected RedisAddr redis:6379, got %s", cfg.RedisAddr)
		}

		if cfg.RedisPassword != "mypassword" {
			t.Errorf("expected RedisPassword mypassword, got %s", cfg.RedisPassword)
		}
	})

	t.Run("short code length", func(t *testing.T) {
		os.Setenv("SHORT_CODE_LENGTH", "8")

		cfg := Load()

		if cfg.ShortCodeLength != 8 {
			t.Errorf("expected ShortCodeLength 8, got %d", cfg.ShortCodeLength)
		}
	})

	t.Run("cache TTL", func(t *testing.T) {
		os.Setenv("CACHE_TTL", "48h")

		cfg := Load()

		if cfg.CacheTTL != 48*time.Hour {
			t.Errorf("expected CacheTTL 48h, got %v", cfg.CacheTTL)
		}
	})
}

func TestGetEnv(t *testing.T) {
	originalValue := os.Getenv("TEST_VAR")
	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_VAR", originalValue)
		} else {
			os.Unsetenv("TEST_VAR")
		}
	}()

	t.Run("existing variable", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test_value")
		result := getEnv("TEST_VAR", "default")
		if result != "test_value" {
			t.Errorf("expected test_value, got %s", result)
		}
	})

	t.Run("missing variable", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		result := getEnv("TEST_VAR", "default")
		if result != "default" {
			t.Errorf("expected default, got %s", result)
		}
	})
}

func TestGetEnvAsInt(t *testing.T) {
	originalValue := os.Getenv("TEST_INT")
	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_INT", originalValue)
		} else {
			os.Unsetenv("TEST_INT")
		}
	}()

	t.Run("valid integer", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		result := getEnvAsInt("TEST_INT", 10)
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("invalid integer", func(t *testing.T) {
		os.Setenv("TEST_INT", "invalid")
		result := getEnvAsInt("TEST_INT", 10)
		if result != 10 {
			t.Errorf("expected default 10, got %d", result)
		}
	})
}

func TestGetEnvAsBool(t *testing.T) {
	originalValue := os.Getenv("TEST_BOOL")
	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_BOOL", originalValue)
		} else {
			os.Unsetenv("TEST_BOOL")
		}
	}()

	t.Run("true values", func(t *testing.T) {
		trueValues := []string{"true", "True", "TRUE", "1", "t", "T"}
		for _, val := range trueValues {
			os.Setenv("TEST_BOOL", val)
			result := getEnvAsBool("TEST_BOOL", false)
			if !result {
				t.Errorf("expected true for value %s, got false", val)
			}
		}
	})

	t.Run("false values", func(t *testing.T) {
		falseValues := []string{"false", "False", "FALSE", "0", "f", "F"}
		for _, val := range falseValues {
			os.Setenv("TEST_BOOL", val)
			result := getEnvAsBool("TEST_BOOL", true)
			if result {
				t.Errorf("expected false for value %s, got true", val)
			}
		}
	})
}

func TestGetEnvAsDuration(t *testing.T) {
	originalValue := os.Getenv("TEST_DURATION")
	defer func() {
		if originalValue != "" {
			os.Setenv("TEST_DURATION", originalValue)
		} else {
			os.Unsetenv("TEST_DURATION")
		}
	}()

	t.Run("valid duration", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "5m")
		result := getEnvAsDuration("TEST_DURATION", 1*time.Minute)
		if result != 5*time.Minute {
			t.Errorf("expected 5m, got %v", result)
		}
	})

	t.Run("invalid duration", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "invalid")
		result := getEnvAsDuration("TEST_DURATION", 1*time.Minute)
		if result != 1*time.Minute {
			t.Errorf("expected default 1m, got %v", result)
		}
	})
}
