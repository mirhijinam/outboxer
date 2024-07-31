package kafka

import (
	"context"

	k "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mirhijinam/outboxer/internal/config"
	"go.uber.org/zap"
)

type Consumer struct {
	c             *k.Consumer
	topic         string
	handleMessage func(ctx context.Context, message []byte) error
	log           *zap.Logger
}

func NewConsumer(kafkaCfg config.KafkaConfig, l *zap.Logger, handleMessage func(ctx context.Context, message []byte) error) (*Consumer, error) {
	c, err := k.NewConsumer(&k.ConfigMap{
		"bootstrap.servers":     kafkaCfg.Brokers,
		"group.id":              kafkaCfg.GroupID,
		"auto.offset.reset":     kafkaCfg.OffsetReset,
		"broker.address.family": "v4",
	})
	if err != nil {
		return nil, err
	}

	err = c.Subscribe(kafkaCfg.Topic, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		c:             c,
		topic:         kafkaCfg.Topic,
		handleMessage: handleMessage,
		log:           l}, nil
}

func (c *Consumer) Close() {
	c.c.Close()
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.log.Info("consuming is stopped. Run: Context was done.")
			return
		default:
			msg, err := c.c.ReadMessage(-1)
			if err == nil {
				if err := c.handleMessage(ctx, msg.Value); err != nil {
					c.log.Error("failed to handle message. StartConsuming:", zap.Error(err))
				}
			} else {
				c.log.Error("failed to consume. StartConsuming:", zap.Error(err), zap.Any("kafka msg", msg))
			}
		}
	}
}
