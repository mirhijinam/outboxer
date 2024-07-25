package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/mirhijinam/outboxer/internal/config"
)

func MustOpenDB(ctx context.Context, cfg config.Config) *sql.DB {
	// construct the dsn
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.DBConfig.PgUser, cfg.DBConfig.PgPassword, cfg.DBConfig.PgHost, cfg.DBConfig.PgPort, cfg.DBConfig.PgDatabase,
	)

	// open db
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("failed to open database", "error", err.Error())
		os.Exit(1)
	}

	// check if db is alive
	if err = db.PingContext(ctx); err != nil {
		slog.Error("failed to ping database", "error", err.Error())
		os.Exit(1)
	}

	return db
}
