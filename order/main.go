package main

import (
	"context"
	"os"
	"strings"

	"order/consumer"
	"order/core"
	"order/jwtauth"
	"order/logger"
	"order/queue/kafka"
	"order/repo/mysql"
	"order/route"
)

func main() {

	dsn := os.Getenv("DB_DSN")
	dbTable := os.Getenv("DB_TABLE")
	kafkaHosts := strings.Split(os.Getenv("KAFKA_HOST"), ",")

	log := logger.New()

	dbRepo, err := mysql.NewRepo(dsn, dbTable, 10)
	if err != nil {
		log.Fatal(err)
	}

	service := core.NewService(
		dbRepo,
		jwtauth.New(),
		kafka.NewProducer(kafkaHosts, log),
	)

	cons := consumer.Port{
		Sub: kafka.NewConsumer(kafkaHosts, log),
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
