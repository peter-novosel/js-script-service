package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/peter-novosel/js-script-service/internal/config"
	"github.com/peter-novosel/js-script-service/internal/db"
	"github.com/peter-novosel/js-script-service/internal/logger"
	"github.com/peter-novosel/js-script-service/internal/router"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize structured logger
	log := logger.Init()

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	if err := db.Init(cfg); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Set up the router
	r := router.Setup(cfg)

	// Start HTTP server
	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("graceful shutdown failed: %v", err)
	} else {
		log.Info("Server stopped gracefully")
	}
}
