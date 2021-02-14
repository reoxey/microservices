package main

import (
	"context"
	"os"
	"strings"

	"shipping/consumer"
	"shipping/jwtauth"
	"shipping/logger"
	"shipping/queue/kafka"
	"shipping/repo/mysql"
	"shipping/route"
	"shipping/core"
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
	)

	cons := consumer.Port{
		Sub: kafka.NewConsumer(
			kafkaHosts,
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
