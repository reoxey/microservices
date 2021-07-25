package mysql

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"cart/config"
	"cart/core"
)

type mysqlRepo struct {
	db    *sql.DB
	table string
}

var (
	rowsEmpty      = sql.ErrNoRows
	NoRowsAffected = errors.New("no_rows_affected")
	noIdUpdate     = errors.New("no_id_for_update")
)

func NewRepo(conf *config.Config) (core.CartRepo, error) {
	db, e := sql.Open("mysql", conf.MysqlDSN)
	if e != nil {
		return nil, e
	}
	e = db.Ping()
	if e != nil {
		return nil, e // proper error handling instead of panic in your app
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(conf.MysqlPoolSize)
	db.SetMaxIdleConns(conf.MysqlPoolSize)

	return &mysqlRepo{
		db:    db,
		table: conf.MysqlTable,
	}, nil
}
