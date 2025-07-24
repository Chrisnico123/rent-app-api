package database

import (
	"context"
	"fmt"
	"rent-application/configs"
	"strconv"
	"time"

	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func NewPostgreDatabase(cfg configs.DBConfig) *pgxpool.Pool {
	// Inisialisasi zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	logger.Info().Str("dsn", dsn).Msg("Database DSN")

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error().Str("dsn", dsn).Err(err).Msg("Failed to parse configuration")
		panic(fmt.Sprintf(`Failed to parse configuration %s`, dsn))
	}

	minConnsInt, err := strconv.Atoi(cfg.ConnectionPool.DB_POOL_MIN)
	if err != nil {
		logger.Error().Str("minimum connections", cfg.ConnectionPool.DB_POOL_MIN).Err(err).Msg("DB_POOL_MIN expected to be integer")
	}
	maxConnsInt, err := strconv.Atoi(cfg.ConnectionPool.DB_POOL_MAX)
	if err != nil {
		logger.Error().Str("maximum connections", cfg.ConnectionPool.DB_POOL_MAX).Err(err).Msg("DB_POOL_MAX expected to be integer")
	}

	poolConfig.MinConns = int32(minConnsInt)
	poolConfig.MaxConns = int32(maxConnsInt)
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Error().Str("dsn", dsn).Err(err).Msg("Failed to apply pool configuration")
	}

	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := pool.Ping(c); err != nil {
		logger.Error().Err(err).Msg("Failed to ping database")
	} else {
		logger.Info().Str("dsn", dsn).Msg("Database connected")
	}

	return pool
}
