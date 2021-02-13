package main

import (
	"context"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"

	"product/cache"
	"product/catalogpb"
	"product/consumer"
	"product/core"
	"product/jwtauth"
	"product/logger"
	"product/queue/kafka"
	"product/repo/mysql"
	"product/route"
)

func main() {

	dsn := os.Getenv("DB_DSN")
	dbTable := os.Getenv("DB_TABLE")
	redisHost := os.Getenv("REDIS")
	kafkaHosts := strings.Split(os.Getenv("KAFKA_HOST"), ",")

	log := logger.New()

	dbRepo, err := mysql.NewRepo(dsn, dbTable, 10)
	if err != nil {
		log.Fatal(err)
	}

	service := core.NewService(
		dbRepo,
		cache.Redis(redisHost),
		jwtauth.New(),
		kafka.NewProducer(kafkaHosts, log),
		)

	go func() {
		l, err := net.Listen("tcp", ":9001")
		if err != nil {
			log.Println(err)
		}

		grpcOb := grpc.NewServer()

		s := catalogpb.NewGrpc(service)

		catalogpb.RegisterCatalogServer(grpcOb, s)

		log.Println("gRPC started :9001")
		if err = grpcOb.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	cons := consumer.Port{
		Sub: kafka.NewConsumer(
			kafkaHosts,
			log,
		),
		Service: service,
		Log:     log,
	}
	go cons.Run(context.Background())

	r := route.New(log, true)

	r.Handle(service)

	if err = r.Run(":8002"); err != nil {
		log.Fatal(err)
	}
}
