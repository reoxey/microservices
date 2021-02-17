package rabbit

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"

	"order/core"
	"order/logger"
)

type queue struct {
	conn *amqp.Connection
	log *logger.Logger
}

func (q queue) Publish(ctx context.Context, msg *core.Message) (err error) {
	ch, err := q.conn.Channel()
	if err != nil {
		q.log.Println("ERROR:rabbit.Publish."+msg.Topic, err)
		return
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		msg.Topic, // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		q.log.Println("ERROR:rabbit.Publish."+msg.Topic, err)
		return
	}

	str, err := json.Marshal(&msg)
	if err != nil {
		q.log.Println("ERROR:rabbit.Publish."+msg.Topic, err)
		return
	}

	err = ch.Publish(
		"direct", // exchange
		"",            // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         str,
		})
	if err != nil {
		q.log.Println("ERROR:rabbit.Publish."+msg.Topic, err)
	}
	return
}

func NewProducer(dsn string, log *logger.Logger) core.Publisher {

	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Println("ERROR:rabbit.NewProducer", err)
		return queue{}
	}

	return queue{
		conn: conn,
	}
}
