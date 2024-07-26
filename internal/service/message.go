package service

import (
	"context"

	"github.com/mirhijinam/outboxer/internal/model"
)

type messageRepository interface {
	Create(ctx context.Context, msg model.Message) error
}

type Service struct {
	messageRepository messageRepository
}

func New(mr messageRepository) *Service {
	return &Service{
		messageRepository: mr,
	}
}

func (s *Service) Create(ctx context.Context, msg model.Message) error {
	return s.messageRepository.Create(ctx, msg)
}
