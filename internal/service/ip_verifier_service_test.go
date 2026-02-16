package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockIPVerifierRepo is a mock implementation of domain.IPVerifierRepo
type MockIPVerifierRepo struct {
	GetCountryByIPFunc func(ctx context.Context, ipAddress string) (string, error)
	HealthCheckFunc    func(ctx context.Context) error
}

func (m *MockIPVerifierRepo) GetCountryByIP(ctx context.Context, ipAddress string) (string, error) {
	if m.GetCountryByIPFunc != nil {
		return m.GetCountryByIPFunc(ctx, ipAddress)
	}
	return "US", nil
}

func (m *MockIPVerifierRepo) HealthCheck(ctx context.Context) error {
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc(ctx)
	}
	return nil
}

func TestVerifyIP_Success_Allowed(t *testing.T) {
	mockRepo := &MockIPVerifierRepo{
		GetCountryByIPFunc: func(ctx context.Context, ipAddress string) (string, error) {
			return "US", nil
		},
	}

	service := NewIPVerifierService(mockRepo)
	ctx := context.Background()

	result, err := service.VerifyIP(ctx, "8.8.8.8", []string{"US", "CA"})

	require.NoError(t, err)
	assert.Equal(t, "8.8.8.8", result.IP)
	assert.Equal(t, "US", result.Country)
	assert.True(t, result.Allowed)
}

func TestVerifyIP_Success_NotAllowed(t *testing.T) {
	mockRepo := &MockIPVerifierRepo{
		GetCountryByIPFunc: func(ctx context.Context, ipAddress string) (string, error) {
			return "CN", nil
		},
	}

	service := NewIPVerifierService(mockRepo)
	ctx := context.Background()

	result, err := service.VerifyIP(ctx, "1.2.3.4", []string{"US", "CA"})

	require.NoError(t, err)
	assert.Equal(t, "1.2.3.4", result.IP)
	assert.Equal(t, "CN", result.Country)
	assert.False(t, result.Allowed)
}

func TestVerifyIP_EmptyAllowedCountries(t *testing.T) {
	mockRepo := &MockIPVerifierRepo{}
	service := NewIPVerifierService(mockRepo)
	ctx := context.Background()

	result, err := service.VerifyIP(ctx, "8.8.8.8", []string{})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "allowed_countries")
}

func TestVerifyIP_RepoError(t *testing.T) {
	mockRepo := &MockIPVerifierRepo{
		GetCountryByIPFunc: func(ctx context.Context, ipAddress string) (string, error) {
			return "", errors.New("invalid IP address")
		},
	}

	service := NewIPVerifierService(mockRepo)
	ctx := context.Background()

	result, err := service.VerifyIP(ctx, "invalid-ip", []string{"US"})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid IP address")
}
