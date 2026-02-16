package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	Environment     string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	GeoIPPath string
}

// Load reads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ReadTimeout:     getDurationEnv("READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 30*time.Second),
			Environment:     getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			GeoIPPath: getEnv("GEOIP_DB_PATH", "data/GeoLite2-Country.mmdb"),
		},
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate ensures all required configuration is present and valid
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	if c.Database.GeoIPPath == "" {
		return fmt.Errorf("GeoIP database path cannot be empty")
	}

	// Validate port is a valid number
	if _, err := strconv.Atoi(c.Server.Port); err != nil {
		return fmt.Errorf("invalid port number: %s", c.Server.Port)
	}

	return nil
}

// GetAddress returns the full server address (e.g., ":8080")
func (c *Config) GetAddress() string {
	return ":" + c.Server.Port
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv retrieves a duration from environment or returns default
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
