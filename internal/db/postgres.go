package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/peter-novosel/js-script-service/internal/config"
	"github.com/peter-novosel/js-script-service/internal/logger"
)

var conn *pgxpool.Pool

// Init initializes the PostgreSQL connection pool
func Init(cfg *config.Config) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	conn, err = pgxpool.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	// Verify connection
	if err := conn.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	logger.Init().Info("Connected to PostgreSQL")
	return nil
}

// MustInit exits on DB init failure
func MustInit() {
	cfg := config.Load()
	if err := Init(cfg); err != nil {
		logger.Init().Fatalf("DB init failed: %v", err)
	}
}

// Close closes the PostgreSQL connection pool
func Close() {
	if conn != nil {
		conn.Close()
	}
}
