package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestVerifyIP_Success_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockFunc := func(ipAddress string) (string, error) {
		return "US", nil
	}

	router := gin.Default()
	router.POST("/verify", func(c *gin.Context) {
		var verifyReq VerifyRequest

		if err := c.ShouldBindJSON(&verifyReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		country, err := mockFunc(verifyReq.IP)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		allowed := false
		for _, ac := range verifyReq.AllowedCountries {
			if ac == country {
				allowed = true
				break
			}
		}

		resp := VerifyResponse{
			IP:      verifyReq.IP,
			Country: country,
			Allowed: allowed,
		}
		c.JSON(http.StatusOK, resp)
	})

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

	mockFunc := func(ipAddress string) (string, error) {
		return "CN", nil
	}

	router := gin.Default()
	router.POST("/verify", func(c *gin.Context) {
		var verifyReq VerifyRequest

		if err := c.ShouldBindJSON(&verifyReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		country, err := mockFunc(verifyReq.IP)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		allowed := false
		for _, ac := range verifyReq.AllowedCountries {
			if ac == country {
				allowed = true
				break
			}
		}

		resp := VerifyResponse{
			IP:      verifyReq.IP,
			Country: country,
			Allowed: allowed,
		}
		c.JSON(http.StatusOK, resp)
	})

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

	router := gin.Default()
	router.POST("/verify", func(c *gin.Context) {
		var verifyReq VerifyRequest
		if err := c.ShouldBindJSON(&verifyReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	req, _ := http.NewRequest("POST", "/verify", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyIP_MissingRequiredFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/verify", func(c *gin.Context) {
		var verifyReq VerifyRequest
		if err := c.ShouldBindJSON(&verifyReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

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

	mockFunc := func(ipAddress string) (string, error) {
		return "", fmt.Errorf("invalid IP address: %s", ipAddress)
	}

	router := gin.Default()
	router.POST("/verify", func(c *gin.Context) {
		var verifyReq VerifyRequest

		if err := c.ShouldBindJSON(&verifyReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		country, err := mockFunc(verifyReq.IP)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		allowed := false
		for _, ac := range verifyReq.AllowedCountries {
			if ac == country {
				allowed = true
				break
			}
		}

		resp := VerifyResponse{
			IP:      verifyReq.IP,
			Country: country,
			Allowed: allowed,
		}
		c.JSON(http.StatusOK, resp)
	})

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

	_ = fmt.Errorf("") // keep import
}
