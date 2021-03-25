package core

import "context"

type Message struct {
	Topic string
	Msg   []byte
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string, rec chan *Message)
}

type Publisher interface {
	Publish(ctx context.Context, msg *Message) error
}
