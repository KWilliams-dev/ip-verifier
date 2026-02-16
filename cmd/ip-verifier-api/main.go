package main

import (
	"context"
	"errors"
	"ip-verifier/internal/api/handler"
	"ip-verifier/internal/config"
	"ip-verifier/internal/repo"
	"ip-verifier/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	slog.Info("Configuration loaded",
		"port", cfg.Server.Port,
		"environment", cfg.Server.Environment,
		"geoip_path", cfg.Database.GeoIPPath,
	)

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Open GeoIP database
	db, err := geoip2.Open(cfg.Database.GeoIPPath)
	if err != nil {
		slog.Error("Failed to open GeoIP database", "error", err, "path", cfg.Database.GeoIPPath)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("GeoIP database opened successfully")

	// Initialize layers
	ipRepo := repo.NewIPVerifierRepo(db)
	ipService := service.NewIPVerifierService(ipRepo)
	slog.Info("Application layers initialized")

	// Setup router
	router := gin.Default()

	router.GET("/api/v1/health", handler.HealthCheck(ipService))
	router.POST("/api/v1/ip-verifier", handler.VerifyIP(ipService))

	// Configure HTTP server
	srv := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("Starting server",
			"address", cfg.GetAddress(),
			"environment", cfg.Server.Environment,
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("Shutdown signal received", "signal", sig.String())

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	} else {
		slog.Info("Server stopped gracefully")
	}
}
