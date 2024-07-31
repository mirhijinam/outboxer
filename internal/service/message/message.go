package message

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mirhijinam/outboxer/internal/model"
)

type messageRepository interface {
	Create(ctx context.Context, msg model.Message) (int, error)
	CreateEvent(ctx context.Context, tx pgx.Tx, payload string) error
	GetEventNew(ctx context.Context) (model.Event, error)
	SetDone(ctx context.Context, id int) error
	GetStats(ctx context.Context) (map[string]int, error)
}

type Service struct {
	messageRepository messageRepository
}

func New(mr messageRepository) *Service {
	return &Service{
		messageRepository: mr,
	}
}

func (s *Service) Create(ctx context.Context, msg model.Message) (int, error) {
	return s.messageRepository.Create(ctx, msg)
}

func (s *Service) GetStats(ctx context.Context) (map[string]int, error) {
	return s.messageRepository.GetStats(ctx)
}
