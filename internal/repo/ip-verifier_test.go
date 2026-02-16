package repo

import (
	"testing"

	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"
)

func TestGetCountryByIP_InvalidIP(t *testing.T) {
	repo := NewIPVerifierRepo(nil)

	tests := []struct {
		name string
		ip   string
	}{
		{"empty string", ""},
		{"invalid format", "not-an-ip"},
		{"incomplete ip", "192.168.1"},
		{"invalid characters", "192.168.1.abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			country, err := repo.GetCountryByIP(tt.ip)
			assert.Error(t, err)
			assert.Empty(t, country)
			assert.Contains(t, err.Error(), "invalid IP address")
		})
	}
}

func TestGetCountryByIP_ValidIP_WithRealDatabase(t *testing.T) {
	// This test requires the actual GeoLite2-Country.mmdb file
	db, err := geoip2.Open("../../data/GeoLite2-Country.mmdb")
	if err != nil {
		t.Skip("Skipping test: GeoLite2-Country.mmdb not found")
		return
	}
	defer db.Close()

	repo := NewIPVerifierRepo(db)

	tests := []struct {
		name            string
		ip              string
		expectedCountry string
	}{
		{"Google DNS US", "8.8.8.8", "US"},
		{"Google DNS IPv6", "2001:4860:4860::8888", "US"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			country, err := repo.GetCountryByIP(tt.ip)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCountry, country)
		})
	}
}

func TestNewIPVerifierRepo(t *testing.T) {
	db, err := geoip2.Open("../../data/GeoLite2-Country.mmdb")
	if err != nil {
		t.Skip("Skipping test: GeoLite2-Country.mmdb not found")
		return
	}
	defer db.Close()

	repo := NewIPVerifierRepo(db)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}
