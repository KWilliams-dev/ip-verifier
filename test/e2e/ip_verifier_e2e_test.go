//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

type VerifyRequest struct {
	IP               string   `json:"ip"`
	AllowedCountries []string `json:"allowed_countries"`
}

type VerifyResponse struct {
	IP      string `json:"ip"`
	Country string `json:"country,omitempty"`
	Allowed bool   `json:"allowed"`
}

func TestMain(m *testing.M) {
	// Wait for server to be ready
	if !waitForServer(baseURL, 10*time.Second) {
		panic("❌ Server not responding. Make sure it's running: go run cmd/ip-verifier-api/main.go")
	}
	m.Run()
}

func waitForServer(url string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url + "/api/v1/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/v1/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"status":"healthy"}`, string(body))
}

func TestVerifyIP_ValidIPInAllowedCountries(t *testing.T) {
	req := VerifyRequest{
		IP:               "8.8.8.8",
		AllowedCountries: []string{"US", "CA"},
	}

	resp := makeVerifyRequest(t, req)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var verifyResp VerifyResponse
	decodeJSON(t, resp, &verifyResp)

	assert.Equal(t, "8.8.8.8", verifyResp.IP)
	assert.Equal(t, "US", verifyResp.Country)
	assert.True(t, verifyResp.Allowed)
}

func TestVerifyIP_ValidIPNotInAllowedCountries(t *testing.T) {
	req := VerifyRequest{
		IP:               "8.8.8.8",
		AllowedCountries: []string{"CN", "RU"},
	}

	resp := makeVerifyRequest(t, req)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var verifyResp VerifyResponse
	decodeJSON(t, resp, &verifyResp)

	assert.Equal(t, "8.8.8.8", verifyResp.IP)
	assert.Equal(t, "US", verifyResp.Country)
	assert.False(t, verifyResp.Allowed)
}

func TestVerifyIP_InvalidIPAddress(t *testing.T) {
	req := VerifyRequest{
		IP:               "invalid-ip",
		AllowedCountries: []string{"US"},
	}

	resp := makeVerifyRequest(t, req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errorResp map[string]string
	decodeJSON(t, resp, &errorResp)
	assert.Contains(t, errorResp["error"], "invalid IP address")
}

func TestVerifyIP_MissingRequiredFields(t *testing.T) {
	reqBody := `{"ip": "8.8.8.8"}`

	resp, err := http.Post(baseURL+"/api/v1/ip-verifier", "application/json", bytes.NewBufferString(reqBody))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestVerifyIP_IPv6Address(t *testing.T) {
	req := VerifyRequest{
		IP:               "2001:4860:4860::8888",
		AllowedCountries: []string{"US"},
	}

	resp := makeVerifyRequest(t, req)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var verifyResp VerifyResponse
	decodeJSON(t, resp, &verifyResp)

	assert.Equal(t, "2001:4860:4860::8888", verifyResp.IP)
	assert.Equal(t, "US", verifyResp.Country)
	assert.True(t, verifyResp.Allowed)
}

// Helper functions

func makeVerifyRequest(t *testing.T, req VerifyRequest) *http.Response {
	body, err := json.Marshal(req)
	require.NoError(t, err)

	resp, err := http.Post(baseURL+"/api/v1/ip-verifier", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)

	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, v interface{}) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = json.Unmarshal(body, v)
	require.NoError(t, err)
}
