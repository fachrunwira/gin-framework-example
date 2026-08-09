package database

import (
	"github.com/fachrunwira/gin-example/lib/env"
)

type dbConfig struct {
	connection string
	host       string
	username   string
	password   string
	dbName     string
	port       string
}

func getConfig() dbConfig {
	return dbConfig{
		connection: env.Get("DATABASE_CONNECTION", "mysql"),
		host:       env.Get("DATABASE_HOST", "localhost"),
		username:   env.Get("DATABASE_USERNAME", "root"),
		password:   env.Get("DATABASE_PASSWORD", ""),
		dbName:     env.Get("DATABASE_NAME", ""),
		port:       env.Get("DATABASE_PORT", "3306"),
	}
}
