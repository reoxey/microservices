package rabbit

import (
	"context"
	"log"

	"github.com/streadway/amqp"

	"shipping/core"
)

func (q queue) Subscribe(ctx context.Context, topic string, rec chan *core.Message) {

	defer q.conn.Close()

	ch, err := q.conn.Channel()
	if err != nil {
		q.log.Println("ERROR:rabbit.Subscribe", err)
		return
	}
	defer ch.Close()

	que, err := ch.QueueDeclare(
		topic, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		q.log.Println("ERROR:rabbit.Subscribe", err)
		return
	}

	err = ch.QueueBind(
		que.Name, // queue name
		"",       // routing key
		"direct", // exchange
		false,
		nil)
	if err != nil {
		q.log.Println("ERROR:rabbit.Subscribe", err)
		return
	}

	msgs, err := ch.Consume(
		que.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		q.log.Println("ERROR:rabbit.Subscribe", err)
		return
	}

	for d := range msgs {
		rec <- &core.Message{
			Topic: topic,
			Msg:   d.Body,
		}
	}
}

func NewConsumer(dsn string) core.Subscriber {

	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Println("ERROR:rabbit.NewConsumer", err)
		return queue{}
	}

	return queue{
		conn: conn,
	}
}
