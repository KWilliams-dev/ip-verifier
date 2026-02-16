package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ip-verifier/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockIPVerifierService is a mock implementation of domain.IPVerifierService
type MockIPVerifierService struct {
	VerifyIPFunc    func(ctx context.Context, ip string, allowedCountries []string) (*domain.VerifyResult, error)
	HealthCheckFunc func(ctx context.Context) error
}

func (m *MockIPVerifierService) VerifyIP(ctx context.Context, ip string, allowedCountries []string) (*domain.VerifyResult, error) {
	if m.VerifyIPFunc != nil {
		return m.VerifyIPFunc(ctx, ip, allowedCountries)
	}
	return nil, nil
}

func (m *MockIPVerifierService) HealthCheck(ctx context.Context) error {
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc(ctx)
	}
	return nil
}

func TestVerifyIP_Success_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockIPVerifierService{
		VerifyIPFunc: func(ctx context.Context, ip string, allowedCountries []string) (*domain.VerifyResult, error) {
			return &domain.VerifyResult{
				IP:      ip,
				Country: "US",
				Allowed: true,
			}, nil
		},
	}

	router := gin.Default()
	router.POST("/verify", VerifyIP(mockService))

	reqBody := VerifyRequest{
		IP:               "8.8.8.8",
		AllowedCountries: []string{"US", "CA"},
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp VerifyResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "8.8.8.8", resp.IP)
	assert.Equal(t, "US", resp.Country)
	assert.True(t, resp.Allowed)
}

func TestVerifyIP_Success_NotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockIPVerifierService{
		VerifyIPFunc: func(ctx context.Context, ip string, allowedCountries []string) (*domain.VerifyResult, error) {
			return &domain.VerifyResult{
				IP:      ip,
				Country: "CN",
				Allowed: false,
			}, nil
		},
	}

	router := gin.Default()
	router.POST("/verify", VerifyIP(mockService))

	reqBody := VerifyRequest{
		IP:               "1.2.3.4",
		AllowedCountries: []string{"US", "CA"},
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp VerifyResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "1.2.3.4", resp.IP)
	assert.Equal(t, "CN", resp.Country)
	assert.False(t, resp.Allowed)
}

func TestVerifyIP_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockIPVerifierService{}
	router := gin.Default()
	router.POST("/verify", VerifyIP(mockService))

	req, _ := http.NewRequest("POST", "/verify", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyIP_MissingRequiredFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockIPVerifierService{}
	router := gin.Default()
	router.POST("/verify", VerifyIP(mockService))

	reqBody := VerifyRequest{
		IP: "8.8.8.8",
		// Missing AllowedCountries
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyIP_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockIPVerifierService{
		VerifyIPFunc: func(ctx context.Context, ip string, allowedCountries []string) (*domain.VerifyResult, error) {
			return nil, fmt.Errorf("invalid IP address: %s", ip)
		},
	}

	router := gin.Default()
	router.POST("/verify", VerifyIP(mockService))

	reqBody := VerifyRequest{
		IP:               "invalid-ip",
		AllowedCountries: []string{"US"},
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
