package main

import (
	"ip-verifier/internal/api/handler"
	"ip-verifier/internal/repo"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	router := gin.Default()

	db, err := geoip2.Open("data/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ipRepo := repo.NewIPVerifierRepo(db)

	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.POST("/api/v1/ip-verifier", handler.VerifyIP(ipRepo))

	router.Run(":8080")
}
