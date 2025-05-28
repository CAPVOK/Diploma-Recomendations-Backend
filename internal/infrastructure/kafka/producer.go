package kafka

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type IKafkaProducer interface {
	Send(ctx context.Context, topic, key string, value interface{}) error
}

type kafkaProducer struct {
	brokers []string
	logger  *zap.Logger

	mu      sync.Mutex
	writers map[string]*kafka.Writer
}

func NewKafkaProducer(brokers []string, logger *zap.Logger) IKafkaProducer {
	return &kafkaProducer{
		brokers: brokers,
		logger:  logger.Named("kafka-producer"),
		writers: make(map[string]*kafka.Writer),
	}
}

func (p *kafkaProducer) getWriter(topic string) *kafka.Writer {
	p.mu.Lock()
	defer p.mu.Unlock()

	w, ok := p.writers[topic]
	if ok {
		return w
	}

	w = &kafka.Writer{
		Addr:      kafka.TCP(p.brokers...),
		Topic:     topic,
		Balancer:  &kafka.LeastBytes{},
		BatchSize: 1,
	}
	p.writers[topic] = w
	return w
}

func (p *kafkaProducer) Send(ctx context.Context, topic, key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		p.logger.Error("failed to marshal message", zap.Error(err))
		return err
	}

	p.logger.Info("kafka: producing message",
		zap.String("topic", topic),
		zap.String("key", key),
		zap.String("value", string(bytes)),
	)

	msg := kafka.Message{
		Key:   []byte(key),
		Value: bytes,
		Time:  time.Now(),
	}

	writer := p.getWriter(topic)
	if err := writer.WriteMessages(ctx, msg); err != nil {
		p.logger.Error("kafka: write message failed",
			zap.Error(err),
			zap.String("topic", topic),
		)
		return err
	}

	p.logger.Debug("kafka: message sent",
		zap.String("topic", topic),
		zap.String("key", key),
	)
	return nil
}
