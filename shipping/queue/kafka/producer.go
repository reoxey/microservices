package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"

	"shipping/core"
	"shipping/logger"
)

func NewProducer(dsn []string, log *logger.Logger) core.Publisher {
	return &queue{dsn: dsn, log: log}
}

func (q queue) Publish(ctx context.Context, msg *core.Message) error {

	w := &kafka.Writer{
		Addr:     kafka.TCP(q.dsn[0]),
		Topic:    msg.Topic,
		Balancer: &kafka.LeastBytes{},
		Logger:   q.log,
	}

	str, err := json.Marshal(&msg)
	if err != nil {
		q.log.Println("ERROR:kafka.Publish."+msg.Topic, err)
		return err
	}

	err = w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("one"),
		Value: str,
	})

	if err != nil {
		q.log.Println("ERROR:kafka.Publish."+msg.Topic, err)
	}

	return err
}
