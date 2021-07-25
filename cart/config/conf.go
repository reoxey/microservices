// Author: reoxey
// Date: 09:34

package config

import (
	"os"
	"strings"
)

type Config struct {
	MysqlDSN      string
	MysqlTable    string
	MysqlPoolSize int
	KafkaHosts    []string
	GrpcPort      string
	HttpPort      string
}

func New() *Config {
	return &Config{
		MysqlDSN:      os.Getenv("DB_DSN"),
		MysqlTable:    os.Getenv("DB_TABLE"),
		MysqlPoolSize: 10,
		KafkaHosts:    strings.Split(os.Getenv("KAFKA_HOST"), ","),
		GrpcPort:      os.Getenv("PRODUCT_GRPC"),
		HttpPort:      ":8003",
	}
}
