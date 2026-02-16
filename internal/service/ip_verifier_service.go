package service

import (
	"context"
	"ip-verifier/internal/domain"
	apperrors "ip-verifier/internal/errors"
)

type ipVerifierService struct {
	repo domain.IPVerifierRepo
}

// NewIPVerifierService creates a new instance of IPVerifierService
func NewIPVerifierService(repo domain.IPVerifierRepo) domain.IPVerifierService {
	return &ipVerifierService{
		repo: repo,
	}
}

// VerifyIP checks if an IP address is from an allowed country
func (s *ipVerifierService) VerifyIP(ctx context.Context, ip string, allowedCountries []string) (*domain.VerifyResult, error) {
	// Validate input
	if len(allowedCountries) == 0 {
		return nil, apperrors.NewValidationError("allowed_countries cannot be empty", nil)
	}

	// Get country for IP address
	country, err := s.repo.GetCountryByIP(ctx, ip)
	if err != nil {
		return nil, err
	}

	// Check if country is in allowed list
	allowed := contains(allowedCountries, country)

	return &domain.VerifyResult{
		IP:      ip,
		Country: country,
		Allowed: allowed,
	}, nil
}

// contains checks if a string slice contains a specific value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// HealthCheck verifies the repository is healthy
func (s *ipVerifierService) HealthCheck(ctx context.Context) error {
	return s.repo.HealthCheck(ctx)
}
