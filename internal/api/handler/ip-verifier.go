package handler

import (
	"ip-verifier/internal/domain"
	apperrors "ip-verifier/internal/errors"
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

func VerifyIP(ipService domain.IPVerifierService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var verifyReq VerifyRequest

		if err := c.ShouldBindJSON(&verifyReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := ipService.VerifyIP(c.Request.Context(), verifyReq.IP, verifyReq.AllowedCountries)
		if err != nil {
			status := apperrors.GetHTTPStatus(err)
			message := apperrors.GetMessage(err)
			c.JSON(status, gin.H{"error": message})
			return
		}

		resp := VerifyResponse{
			IP:      result.IP,
			Country: result.Country,
			Allowed: result.Allowed,
		}
		c.JSON(http.StatusOK, resp)
	}
}
