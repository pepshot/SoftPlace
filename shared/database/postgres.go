package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pepshot/SoftPlace/shared/config"
	"github.com/pepshot/SoftPlace/shared/logger"
)

func NewPostgresPool(ctx context.Context, cfg config.DatabaseConfig, logger *logger.Logger) (*pgxpool.Pool, error) {
	dsn := buildDSN(cfg)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("failed to parse postgres config", "error", err)
		return nil, err
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Error("failed to create postgres pool", "error", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		logger.Error("failed to ping postgres", "error", err)
		return nil, err
	}

	logger.Info(
		"postgres connection established",
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.Name,
		"maxConns", cfg.MaxConns,
		"minConns", cfg.MinConns,
	)

	return pool, nil
}

func buildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
}
