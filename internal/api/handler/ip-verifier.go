package handler

import (
	"ip-verifier/internal/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VerifyRequest struct {
	IP               string   `json:"ip" binding:"required"`
	AllowedCountries []string `json:"allowed_countries" binding:"required"`
}

type VerifyResponse struct {
	IP      string `json:"ip"`
	Country string `json:"country,omitempty"`
	Allowed bool   `json:"allowed"`
}

func VerifyIP(ipRepo *repo.IPVerifierRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var verifyReq VerifyRequest

		if err := c.ShouldBindJSON(&verifyReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		country, err := ipRepo.GetCountryByIP(verifyReq.IP)
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
	}
}
