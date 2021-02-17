package main

import (
	"context"

	"order/consumer"
	"order/core"
	"order/jwtauth"
	"order/logger"
	"order/queue/kafka"
	"order/repo/mysql"
	"order/route"
)

func main() {

	log := logger.New()

	dsn := ""

	dbRepo, err := mysql.NewRepo(dsn, "orders", 10)
	if err != nil {
		log.Fatal(err)
	}

	service := core.NewService(
		dbRepo,
		jwtauth.New(),
		kafka.NewProducer([]string{
			"localhost:9092",
		}, log),
	)

	cons := consumer.Port{
		Sub: kafka.NewConsumer(
			[]string{"localhost:9092"},
			log,
		),
		Service:   service,
		Log:       log,
	}
	go cons.Run(context.Background())

	r := route.New(log, true)

	r.Handle(service)

	if err = r.Run(":8005"); err != nil {
		log.Fatal(err)
	}
}
