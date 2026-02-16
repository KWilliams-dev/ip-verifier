package handler

import (
	"ip-verifier/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// HealthCheck creates a handler for health check endpoint
func HealthCheck(ipService domain.IPVerifierService) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ipService.HealthCheck(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, HealthResponse{
				Status:  "unhealthy",
				Message: "GeoIP database unavailable",
			})
			return
		}

		c.JSON(http.StatusOK, HealthResponse{
			Status: "healthy",
		})
	}
}
