// Author: reoxey
// Date: 28-07-2021 09:18

// wire.go

// +build wireinject

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/wire"
	"google.golang.org/grpc"

	"payment/catalogpb"
	"payment/config"
	"payment/consumer"
	"payment/core"
	"payment/jwtauth"
	"payment/queue/kafka"
	"payment/repo/mysql"
	"payment/route"
)

type App struct {
	router http.Handler
	port   *consumer.Port
}

func main() {

	conf := config.New()

	app, err := Init(conf)
	if err != nil {
		log.Fatal(err)
	}

	go app.port.Run(context.Background())

	srv := &http.Server{
		Addr:    conf.HttpPort,
		Handler: app.router,
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

func NewCatalogConn() grpc.ClientConnInterface {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Error(err)
	}
	defer conn.Close()

	return conn
}

func Init(conf *config.Config) (*App, error) {
	wire.Build(jwtauth.New, mysql.NewRepo, NewCatalogConn, catalogpb.NewCatalogClient,
		kafka.NewProducer, core.NewService, kafka.NewConsumer, consumer.InitRun, route.New,
		wire.Struct(new(App), "*"))
	return nil, nil
}

// func Init() core.CartService {
// 	dbRepo, err := mysql.NewRepo(dsn, dbTable, 10)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	conn, err := grpc.Dial(grpcPort, grpc.WithInsecure())
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	defer conn.Close()
//
// 	return core.NewService(
// 		dbRepo,
// 		jwtauth.New(),
// 		catalogpb.NewCatalogClient(conn),
// 		kafka.NewProducer(kafkaHosts),
// 	)
// }
