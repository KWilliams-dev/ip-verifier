package repo

import (
	"context"
	"ip-verifier/internal/domain"
	apperrors "ip-verifier/internal/errors"
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
		return "", apperrors.NewValidationError("Invalid IP address", nil)
	}

	record, err := r.db.Country(ip)
	if err != nil {
		return "", apperrors.NewInternalError("Failed to lookup IP address", err)
	}

	return record.Country.IsoCode, nil
}

// HealthCheck verifies the GeoIP database is accessible
func (r *IPVerifierRepo) HealthCheck(ctx context.Context) error {
	if r.db == nil {
		return apperrors.NewInternalError("GeoIP database not initialized", nil)
	}
	// Try a simple lookup to verify DB is working
	_, err := r.db.Country(net.ParseIP("8.8.8.8"))
	if err != nil {
		return apperrors.NewInternalError("GeoIP database health check failed", err)
	}
	return nil
}
