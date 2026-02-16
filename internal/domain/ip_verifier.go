package domain

import "context"

// IPVerifierRepo defines the interface for IP geolocation data access
type IPVerifierRepo interface {
	GetCountryByIP(ctx context.Context, ipAddress string) (string, error)
	HealthCheck(ctx context.Context) error
}

// IPVerifierService defines the interface for IP verification business logic
type IPVerifierService interface {
	VerifyIP(ctx context.Context, ip string, allowedCountries []string) (*VerifyResult, error)
}

// VerifyResult represents the result of an IP verification
type VerifyResult struct {
	IP      string
	Country string
	Allowed bool
}
