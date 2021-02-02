package main

import (
	"net"

	"google.golang.org/grpc"

	"product/cache"
	"product/catalog"
	"product/catalogpb"
	"product/jwtauth"
	"product/logger"
	"product/repo/mysql"
	"product/route"
)

func main() {

	dsn := "micro:micro@tcp(127.0.0.1:3306)/micro"

	log := logger.New()

	dbRepo, err := mysql.NewRepo(dsn, "products", 10)
	if err != nil {
		log.Fatal(err)
	}

	service := catalog.NewService(
		dbRepo,
		cache.Redis("localhost"),
		jwtauth.New(),
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

	r := route.New(log, true)

	r.Handle(service)

	if err = r.Run(":8002"); err != nil {
		log.Fatal(err)
	}
}
