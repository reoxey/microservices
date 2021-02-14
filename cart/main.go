package main

import (
	"context"
	"os"
	"strings"

	"google.golang.org/grpc"

	"cart/core"
	"cart/catalogpb"
	"cart/consumer"
	"cart/jwtauth"
	"cart/logger"
	"cart/queue/kafka"
	"cart/repo/mysql"
	"cart/route"
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

	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	service := core.NewService(
		dbRepo,
		jwtauth.New(),
		catalogpb.NewCatalogClient(conn),
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

	if err = r.Run(":8003"); err != nil {
		log.Fatal(err)
	}
}
