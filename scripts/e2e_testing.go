package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type VerifyRequest struct {
	IP               string   `json:"ip"`
	AllowedCountries []string `json:"allowed_countries"`
}

type VerifyResponse struct {
	IP      string `json:"ip"`
	Country string `json:"country,omitempty"`
	Allowed bool   `json:"allowed"`
}

func main() {
	baseURL := "http://localhost:8080"

	// Wait for server to be ready
	fmt.Println("Waiting for server to be ready...")
	if !waitForServer(baseURL, 10*time.Second) {
		fmt.Println("❌ Server not responding. Make sure it's running: go run cmd/ip-verifier-api/main.go")
		os.Exit(1)
	}
	fmt.Println("✅ Server is ready\n")

	// Test 1: Health check
	fmt.Println("=== Test 1: Health Check ===")
	testHealthCheck(baseURL)

	// Test 2: Valid IP in allowed countries
	fmt.Println("\n=== Test 2: Valid IP in Allowed Countries ===")
	testValidIPAllowed(baseURL)

	// Test 3: Valid IP not in allowed countries
	fmt.Println("\n=== Test 3: Valid IP Not in Allowed Countries ===")
	testValidIPNotAllowed(baseURL)

	// Test 4: Invalid IP address
	fmt.Println("\n=== Test 4: Invalid IP Address ===")
	testInvalidIP(baseURL)

	// Test 5: Missing required fields
	fmt.Println("\n=== Test 5: Missing Required Fields ===")
	testMissingFields(baseURL)

	// Test 6: IPv6 address
	fmt.Println("\n=== Test 6: IPv6 Address ===")
	testIPv6(baseURL)

	fmt.Println("\n✅ All tests completed!")
}

func waitForServer(baseURL string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL + "/api/v1/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func testHealthCheck(baseURL string) {
	resp, err := http.Get(baseURL + "/api/v1/health")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))

	if resp.StatusCode == 200 {
		fmt.Println("✅ Health check passed")
	} else {
		fmt.Println("❌ Health check failed")
	}
}

func testValidIPAllowed(baseURL string) {
	req := VerifyRequest{
		IP:               "8.8.8.8",
		AllowedCountries: []string{"US", "CA"},
	}
	resp := makeVerifyRequest(baseURL, req)

	if resp != nil && resp.Allowed && resp.Country == "US" {
		fmt.Println("✅ Test passed: IP allowed in US")
	} else {
		fmt.Println("❌ Test failed")
	}
}

func testValidIPNotAllowed(baseURL string) {
	req := VerifyRequest{
		IP:               "8.8.8.8",
		AllowedCountries: []string{"CN", "RU"},
	}
	resp := makeVerifyRequest(baseURL, req)

	if resp != nil && !resp.Allowed && resp.Country == "US" {
		fmt.Println("✅ Test passed: IP not allowed")
	} else {
		fmt.Println("❌ Test failed")
	}
}

func testInvalidIP(baseURL string) {
	req := VerifyRequest{
		IP:               "invalid-ip",
		AllowedCountries: []string{"US"},
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/api/v1/ip-verifier", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(responseBody))

	if resp.StatusCode == 400 {
		fmt.Println("✅ Test passed: Invalid IP rejected")
	} else {
		fmt.Println("❌ Test failed: Expected 400 status")
	}
}

func testMissingFields(baseURL string) {
	reqBody := `{"ip": "8.8.8.8"}`

	resp, err := http.Post(baseURL+"/api/v1/ip-verifier", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(responseBody))

	if resp.StatusCode == 400 {
		fmt.Println("✅ Test passed: Missing fields rejected")
	} else {
		fmt.Println("❌ Test failed: Expected 400 status")
	}
}

func testIPv6(baseURL string) {
	req := VerifyRequest{
		IP:               "2001:4860:4860::8888",
		AllowedCountries: []string{"US"},
	}
	resp := makeVerifyRequest(baseURL, req)

	if resp != nil && resp.Allowed && resp.Country == "US" {
		fmt.Println("✅ Test passed: IPv6 address handled correctly")
	} else {
		fmt.Println("❌ Test failed")
	}
}

func makeVerifyRequest(baseURL string, req VerifyRequest) *VerifyResponse {
	body, _ := json.Marshal(req)
	fmt.Printf("Request: %s\n", string(body))

	resp, err := http.Post(baseURL+"/api/v1/ip-verifier", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		return nil
	}

	var verifyResp VerifyResponse
	json.Unmarshal(responseBody, &verifyResp)
	return &verifyResp
}
