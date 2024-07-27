package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mirhijinam/outboxer/internal/model"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Repository {
	return &Repository{
		pool: p,
	}
}

func (r *Repository) Create(ctx context.Context, msg model.Message) error {
	query := `INSERT INTO message (content, created_at)
			  VALUES ($1, $2)`

	_, err := r.pool.Exec(ctx, query, msg.Content, msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to exec query: Create: %w", err)
	}

	return nil
}
