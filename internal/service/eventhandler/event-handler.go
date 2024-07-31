package eventhandler

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mirhijinam/outboxer/internal/config"
	"github.com/mirhijinam/outboxer/internal/model"
	"github.com/mirhijinam/outboxer/internal/service/kafka"
	"go.uber.org/zap"
)

type messageRepository interface {
	Create(ctx context.Context, msg model.Message) (int, error)
	CreateEvent(ctx context.Context, tx pgx.Tx, payload string) error
	GetEventNew(ctx context.Context) (model.Event, error)
	SetDone(ctx context.Context, id int) error
}

type eventHandler struct {
	messageRepository messageRepository
	kProducer         *kafka.Producer
	log               *zap.Logger
	cooldown          time.Duration
}

func New(ehConfig config.EventHandlerConfig, mr messageRepository, kp *kafka.Producer, l *zap.Logger) *eventHandler {

	handlerCooldown := time.Duration(ehConfig.CooldownSec) * time.Second
	return &eventHandler{
		messageRepository: mr,
		kProducer:         kp,
		log:               l,
		cooldown:          handlerCooldown,
	}
}

func (eh *eventHandler) Close() {
	if eh.kProducer != nil {
		eh.kProducer.Close()
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
			}

			ev, err := eh.messageRepository.GetEventNew(ctx)
			if err != nil {
				eh.log.Error("failed to get event w/ status 'new'. StartHandlingEvents:", zap.Error(err))
				continue
			}

			if ev.ID == 0 {
				continue
			}

			if err := eh.kProducer.SendMessage(ctx, []byte(ev.Payload)); err != nil {
				eh.log.Error("failed to send message to Kafka. StartHandlingEvents:", zap.Error(err))
				continue
			}

			if err := eh.messageRepository.SetDone(ctx, ev.ID); err != nil {
				eh.log.Error("failed to set event status 'done':", zap.Error(err))
				continue
			}
		}
	}()
}
