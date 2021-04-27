package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	srv := &http.Server{
		Addr:    ":8005",
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
