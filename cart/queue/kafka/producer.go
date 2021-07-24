package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"

	"cart/core"
	log "cart/logger"
)

func NewProducer(dsn []string) core.Publisher {
	return &queue{dsn: dsn}
}

func (q queue) Publish(ctx context.Context, msg *core.Message) error {

	w := &kafka.Writer{
		Addr:     kafka.TCP(q.dsn[0]),
		Topic:    msg.Topic,
		Balancer: &kafka.LeastBytes{},
		Logger:   log.Obj(),
	}

	str, err := json.Marshal(&msg)
	if err != nil {
		return log.Err(err)
	}

	err = w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("one"),
		Value: str,
	})

	if err != nil {
		return log.Err(err)
	}

	return nil
}
