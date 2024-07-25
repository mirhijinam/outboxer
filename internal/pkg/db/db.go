package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mirhijinam/outboxer/internal/config"
)

func MustOpenDB(ctx context.Context, cfg config.Config) (*sql.DB, error) {
	// construct the dsn
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.DBConfig.PgUser, cfg.DBConfig.PgPassword, cfg.DBConfig.PgHost, cfg.DBConfig.PgPort, cfg.DBConfig.PgDatabase,
	)

	// open db
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		// TODO: add a log msg
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// check if db is alive
	if err = db.PingContext(ctx); err != nil {
		// TODO: add a log msg
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
