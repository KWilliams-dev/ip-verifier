package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear any environment variables
	os.Clearenv()

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, 10*time.Second, config.Server.ReadTimeout)
	assert.Equal(t, 10*time.Second, config.Server.WriteTimeout)
	assert.Equal(t, 30*time.Second, config.Server.ShutdownTimeout)
	assert.Equal(t, "development", config.Server.Environment)
	assert.Equal(t, "data/GeoLite2-Country.mmdb", config.Database.GeoIPPath)
}

func TestLoad_CustomValues(t *testing.T) {
	os.Setenv("PORT", "9000")
	os.Setenv("READ_TIMEOUT", "5s")
	os.Setenv("WRITE_TIMEOUT", "5s")
	os.Setenv("SHUTDOWN_TIMEOUT", "15s")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("GEOIP_DB_PATH", "/custom/path/GeoLite2.mmdb")
	defer os.Clearenv()

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "9000", config.Server.Port)
	assert.Equal(t, 5*time.Second, config.Server.ReadTimeout)
	assert.Equal(t, 5*time.Second, config.Server.WriteTimeout)
	assert.Equal(t, 15*time.Second, config.Server.ShutdownTimeout)
	assert.Equal(t, "production", config.Server.Environment)
	assert.Equal(t, "/custom/path/GeoLite2.mmdb", config.Database.GeoIPPath)
}

func TestValidate_InvalidPort(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Port: "invalid",
		},
		Database: DatabaseConfig{
			GeoIPPath: "data/GeoLite2-Country.mmdb",
		},
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid port number")
}

func TestValidate_EmptyPort(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Port: "",
		},
		Database: DatabaseConfig{
			GeoIPPath: "data/GeoLite2-Country.mmdb",
		},
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "port cannot be empty")
}

func TestValidate_EmptyGeoIPPath(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		Database: DatabaseConfig{
			GeoIPPath: "",
		},
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GeoIP database path cannot be empty")
}

func TestGetAddress(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Port: "8080",
		},
	}

	assert.Equal(t, ":8080", config.GetAddress())
}

func TestIsProduction(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		expected    bool
	}{
		{"production environment", "production", true},
		{"development environment", "development", false},
		{"staging environment", "staging", false},
		{"empty environment", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Server: ServerConfig{
					Environment: tt.environment,
				},
			}
			assert.Equal(t, tt.expected, config.IsProduction())
		})
	}
}
