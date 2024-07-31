package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mirhijinam/outboxer/internal/model"
)

type Repository struct {
	pool *pgxpool.Pool
}

type event struct {
	Id            int       `pool:"id"`
	Payload       string    `pool:"payload"`
	CreatedAt     time.Time `pool:"created_at"`
	ReservedUntil time.Time `pool:"reserved_until"`
}

func New(p *pgxpool.Pool) *Repository {
	return &Repository{
		pool: p,
	}
}

func (r *Repository) Create(ctx context.Context, msg model.Message) (lastInsertedId int, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to start transaction. Create: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("rollback failed: %w, original error: %v", rollbackErr, err)
			}
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("commit failed: %w", commitErr)
		}
	}()

	query := `INSERT INTO message (content)
              VALUES ($1)
              RETURNING id`

	var lastInsertId int
	err = tx.QueryRow(ctx, query, msg.Content).Scan(&lastInsertId)
	if err != nil {
		return -1, fmt.Errorf("failed to exec query. Create: %w", err)
	}

	// save the message in the outbox-table
	eventPayload := fmt.Sprintf(`{"id": %d, "content": "%s"}`, lastInsertId, msg.Content)
	if err := r.CreateEvent(ctx, tx, eventPayload); err != nil {
		return -1, fmt.Errorf("failed to create event. Create: %w", err)
	}

	// Commit the transaction after all operations are done
	return lastInsertId, nil
}

func (r *Repository) CreateEvent(ctx context.Context, tx pgx.Tx, payload string) error {
	query := `INSERT INTO event (payload)
			  VALUES ($1)`

	_, err := tx.Exec(ctx, query, payload)
	if err != nil {
		return fmt.Errorf("failed to exec query. CreateEvent: %w", err)
	}

	return nil
}

func (r *Repository) GetEventNew(ctx context.Context) (model.Event, error) {
	query := `SELECT id, payload, created_at, reserved_until FROM event
			  WHERE status = 'new'
			  LIMIT 1`
	row := r.pool.QueryRow(ctx, query)

	var evnt event
	err := row.Scan(&evnt.Id, &evnt.Payload, &evnt.CreatedAt, &evnt.ReservedUntil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Event{}, nil
		}

		return model.Event{}, fmt.Errorf("failed to scan event w/ status 'new'. GetEventNew: %w", err)
	}

	if !evnt.CreatedAt.Equal(evnt.ReservedUntil) {
		return model.Event{}, fmt.Errorf("failed to scan event w/ status 'new'. GetEventNew: %w", errors.New("the event has been already reserved"))
	}

	reservedUntil := evnt.ReservedUntil.Add(10 * time.Minute) // TODO: add it to config
	return model.Event{
		ID:            evnt.Id,
		Payload:       evnt.Payload,
		ReservedUntil: reservedUntil,
	}, nil
}

func (r *Repository) SetDone(ctx context.Context, id int) error {
	query := `UPDATE event
			  SET status = 'done' WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to exec query. SetDone: %w,", err)
	}

	return nil
}

func (r *Repository) GetStats(ctx context.Context) (map[string]int, error) {
	ret := make(map[string]int)
	queries := []string{
		`SELECT COUNT(*) FROM event WHERE status = 'new'`,
		`SELECT COUNT(*) FROM event WHERE status = 'done'`,
	}

	// Выполнение запросов
	for i, query := range queries {
		var count int
		row := r.pool.QueryRow(ctx, query)
		err := row.Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan the row. GetStats: %w", err)
		}

		switch i {
		case 0:
			ret["new"] = count
		case 1:
			ret["done"] = count
		}

		ret["all"] = ret["new"] + ret["done"]
	}

	return ret, nil
}
