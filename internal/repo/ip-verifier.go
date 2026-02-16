package repo

import (
	"context"
	"fmt"
	"ip-verifier/internal/domain"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type IPVerifierRepo struct {
	db *geoip2.Reader
}

// NewIPVerifierRepo creates a new IPVerifierRepo that implements domain.IPVerifierRepo
func NewIPVerifierRepo(db *geoip2.Reader) domain.IPVerifierRepo {
	return &IPVerifierRepo{
		db: db,
	}
}

// GetCountryByIP retrieves the country code for a given IP address
func (r *IPVerifierRepo) GetCountryByIP(ctx context.Context, ipAddress string) (string, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address: %s", ipAddress)
	}

	record, err := r.db.Country(ip)
	if err != nil {
		return "", err
	}

	return record.Country.IsoCode, nil
}

// HealthCheck verifies the GeoIP database is accessible
func (r *IPVerifierRepo) HealthCheck(ctx context.Context) error {
	// Try a simple lookup to verify DB is working
	_, err := r.db.Country(net.ParseIP("8.8.8.8"))
	return err
}
