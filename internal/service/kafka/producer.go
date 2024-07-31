package kafka

import (
	"context"

	k "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mirhijinam/outboxer/internal/config"
	"go.uber.org/zap"
)

type Producer struct {
	p     *k.Producer
	topic string
	log   *zap.Logger
}

func NewProducer(kafkaCfg config.KafkaConfig, l *zap.Logger) (*Producer, error) {
	p, err := k.NewProducer(&k.ConfigMap{"bootstrap.servers": kafkaCfg.Brokers})
	if err != nil {
		return nil, err
	}

	return &Producer{
		p:     p,
		topic: kafkaCfg.Topic,
		log:   l}, nil
}

func (p *Producer) Close() {
	p.p.Close()
}

func (p *Producer) SendMessage(ctx context.Context, message []byte) error {
	deliverych := make(chan k.Event)
	defer close(deliverych)

	err := p.p.Produce(&k.Message{
		TopicPartition: k.TopicPartition{Topic: &p.topic, Partition: k.PartitionAny},
		Value:          message,
	}, deliverych)
	if err != nil {
		return err
	}

	msg := (<-deliverych).(*k.Message)
	if msg.TopicPartition.Error != nil {
		return msg.TopicPartition.Error
	}

	p.log.Info("message delivered:",
		zap.String("topic", *msg.TopicPartition.Topic), zap.Int32("partition", msg.TopicPartition.Partition))

	return nil
}
