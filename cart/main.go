package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"cart/catalogpb"
	"cart/consumer"
	"cart/core"
	"cart/jwtauth"
	"cart/logger"
	"cart/queue/kafka"
	"cart/repo/mysql"
	"cart/route"
)

var dsn = os.Getenv("DB_DSN")
var dbTable = os.Getenv("DB_TABLE")
var grpcPort = os.Getenv("PRODUCT_GRPC")
var kafkaHosts = strings.Split(os.Getenv("KAFKA_HOST"), ",")

func main() {

	service := Init()

	cons := consumer.Port{
		Sub:     kafka.NewConsumer(kafkaHosts),
		Service: service,
	}
	go cons.Run(context.Background())

	r := route.New(true)

	r.Handle(service)

	srv := &http.Server{
		Addr:    ":8003",
		Handler: r,
	}

	go func() { // started http server
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// channel listening for interrupts to ensure graceful shutdown of the http server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Msg("Server exiting")
}

func Init() core.CartService {
	dbRepo, err := mysql.NewRepo(dsn, dbTable, 10)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(grpcPort, grpc.WithInsecure())
	if err != nil {
		log.Error(err)
	}
	defer conn.Close()

	return core.NewService(
		dbRepo,
		jwtauth.New(),
		catalogpb.NewCatalogClient(conn),
		kafka.NewProducer(kafkaHosts),
	)
}
