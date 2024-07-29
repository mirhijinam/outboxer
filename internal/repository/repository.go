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
	Id          int           `pool:"id"`
	Payload     string        `pool:"payload"`
	CreatedAt   time.Time     `pool:"created_at"`
	ReservedFor time.Duration `pool:"reserved_for"`
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
				err = errors.Join(rollbackErr)
			}
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = errors.Join(commitErr)
		}
	}()

	query := `INSERT INTO message (content)
			  VALUES ($1)`

	lastInsertId := 0
	err = tx.QueryRow(ctx, query, msg.Content).Scan(&lastInsertId)
	if err != nil {
		return -1, fmt.Errorf("failed to exec query. Create: %w", err)
	}

	// save the message in the outbox-table
	eventPayload := fmt.Sprintf(`{"id": %d, "content": "%s"}`, lastInsertId, msg.Content)
	reservedFor := 5 * time.Hour // TODO: add it to config
	if err := r.CreateEvent(ctx, tx, eventPayload, reservedFor); err != nil {
		return -1, fmt.Errorf("failed to create event. Create: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to commit transaction. Create: %w", err)
	}
	return lastInsertId, nil
}

func (r *Repository) CreateEvent(ctx context.Context, tx pgx.Tx, payload string, reservedFor time.Duration) error {
	query := `INSERT INTO event (payload, reserved_for)
			  VALUES ($1, $2)`

	_, err := tx.Exec(ctx, query, payload, reservedFor)
	if err != nil {
		return fmt.Errorf("failed to exec query. CreateEvent: %w", err)
	}

	return nil

}

func (r *Repository) GetEventNew(ctx context.Context) (model.Event, error) {
	query := `SELECT id, payload, reserved_for FROM event
			  WHERE status = 'new'
			  LIMIT 1`
	row := r.pool.QueryRow(ctx, query)

	var evnt event
	err := row.Scan(&evnt.Id, &evnt.Payload, &evnt.CreatedAt, &evnt.ReservedFor)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Event{}, nil
		}

		return model.Event{}, fmt.Errorf("failed to scan event w/ status 'new'. GetEventNew: %w", err)
	}

	if time.Now().Before(evnt.CreatedAt.Add(evnt.ReservedFor)) {
		return model.Event{}, fmt.Errorf("failed to scan event w/ status 'new'. GetEventNew: %w", errors.New("the event has been already reserved"))
	}
	return model.Event{
		ID:      evnt.Id,
		Payload: evnt.Payload,
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
