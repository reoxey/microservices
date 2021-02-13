package main

import (
	"context"

	"shipping/consumer"
	"shipping/jwtauth"
	"shipping/logger"
	"shipping/queue/kafka"
	"shipping/repo/mysql"
	"shipping/route"
	"shipping/core"
)

func main() {

	dsn := "micro:micro@tcp(127.0.0.1:3306)/micro"

	log := logger.New()

	dbRepo, err := mysql.NewRepo(dsn, "shipping", 10)
	if err != nil {
		log.Fatal(err)
	}

	service := core.NewService(
		dbRepo,
		jwtauth.New(),
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

	if err = r.Run(":8004"); err != nil {
		log.Fatal(err)
	}

}
