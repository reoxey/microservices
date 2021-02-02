package main

import (
	"user/cache"
	"user/jwtauth"
	"user/logger"
	"user/profile"
	"user/repo/mysql"
	"user/route"
)

func main() {

	dsn := "micro:micro@tcp(127.0.0.1:3306)/micro"

	log := logger.New()

	dbRepo, err := mysql.NewRepo(dsn, "users", 10)
	if err != nil {
		log.Fatal(err)
	}

	service := profile.NewService(
		dbRepo,
		cache.Redis("localhost"),
		jwtauth.New(),
		)

	r := route.New(log, true)

	r.Handle(service)

	if err = r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
