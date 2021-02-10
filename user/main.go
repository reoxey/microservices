package main

import (
	"os"

	"user/cache"
	"user/jwtauth"
	"user/logger"
	"user/core"
	"user/repo/mysql"
	"user/route"
)

func main() {

	dsn := os.Getenv("DB_DSN")
	dbTable := os.Getenv("DB_TABLE")
	redisHost := os.Getenv("REDIS")

	log := logger.New()

	dbRepo, err := mysql.NewRepo(dsn, dbTable, 10)
	if err != nil {
		log.Fatal(err)
	}

	service := core.NewService(
		dbRepo,
		cache.Redis(redisHost),
		jwtauth.New(),
		)

	r := route.New(log, true)

	r.Handle(service)

	if err = r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
