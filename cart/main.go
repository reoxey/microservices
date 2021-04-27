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
	grpcPort := os.Getenv("PRODUCT_GRPC")
	kafkaHosts := strings.Split(os.Getenv("KAFKA_HOST"), ",")

	log := logger.New()

	dbRepo, err := mysql.NewRepo(dsn, dbTable, 10)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(grpcPort, grpc.WithInsecure())
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

	srv := &http.Server{
		Addr:    ":8003",
		Handler: r,
	}

	go func() { // started http server
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// channel listening for interrupts to ensure graceful shutdown of the http server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
