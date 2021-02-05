package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"

	"product/core"
	"product/logger"
)

type queue struct {
	dsn []string
	log *logger.Logger
}

func (q queue) Subscribe(ctx context.Context, topic string, rec chan *core.Message) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     q.dsn,
		Topic:       topic,
		GroupID:     "main",
		StartOffset: kafka.FirstOffset,
		MinBytes:    5,
		MaxBytes:    1e6,
		MaxWait:     1 * time.Second,
		Logger:      q.log,
	})
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			q.log.Println("ERROR:kafka.Subscribe", err)
			continue
		}
		rec <- &core.Message{
			Topic: m.Topic,
			Msg:   m.Value,
		}
	}
}

func NewConsumer(dsn []string, log *logger.Logger) core.Subscriber {
	return queue{
		dsn: dsn,
		log: log,
	}
}