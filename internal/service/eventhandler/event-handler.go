package eventhandler

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mirhijinam/outboxer/internal/config"
	"github.com/mirhijinam/outboxer/internal/model"
	"go.uber.org/zap"
)

type messageRepository interface {
	Create(ctx context.Context, msg model.Message) (int, error)
	CreateEvent(ctx context.Context, tx pgx.Tx, payload string, reservedFor time.Duration) error
	GetEventNew(ctx context.Context) (model.Event, error)
	SetDone(ctx context.Context, id int) error
}

type eventHandler struct {
	messageRepository messageRepository
	log               *zap.Logger
	cooldown          time.Duration
}

func New(ehConfig config.EventHandlerConfig, mr messageRepository, l *zap.Logger) *eventHandler {

	handlerCooldown := time.Duration(ehConfig.CooldownSec) * time.Second
	return &eventHandler{
		messageRepository: mr,
		log:               l,
		cooldown:          handlerCooldown,
	}
}

func (eh *eventHandler) StartHandlingEvents(ctx context.Context) {

	ticker := time.NewTicker(
		eh.cooldown,
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				eh.log.Info("handling is stopped. StartHandlingEvents: Context was done.")
				return
			case <-ticker.C:
				// waiting
			}

			ev, err := eh.messageRepository.GetEventNew(ctx)
			if err != nil {
				eh.log.Error("failed to get event w/ status 'new'. StartHandlingEvents: %w", zap.Error(err))
				continue
			}

			if ev.ID == 0 {
				continue
			}

			// TODO: implement eh.SendMessage

			if err := eh.messageRepository.SetDone(ctx, ev.ID); err != nil {
				eh.log.Error("failed to set event status 'done': %w", zap.Error(err))
				continue
			}
		}
	}()
}
