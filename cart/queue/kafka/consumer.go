package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"

	"cart/core"
	"cart/logger"
)

type queue struct {
	dsn []string
	log *logger.Logger
}

func (q queue) Subscribe(ctx context.Context, topic string, rec chan *core.Message) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: q.dsn,
		Topic:   topic,
		GroupID: "cart",
		//StartOffset: kafka.LastOffset,
		MinBytes: 5,
		MaxBytes: 1e6,
		MaxWait: 1 * time.Second,
		//Logger: q.log,
	})
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			q.log.Println("ERROR:kafka.Subscribe", err)
			continue
		}

		var msg *core.Message
		if err = json.Unmarshal(m.Value, &msg); err != nil {
			continue
		}

		rec <- msg
	}
}

func NewConsumer(dsn []string, log *logger.Logger) core.Subscriber {
	return queue{
		dsn: dsn,
		log: log,
	}
}
