package database

import "github.com/fachrunwira/ebookamd-api/lib/env"

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
		connection: env.GetEnv("DATABASE_CONNECTION", "mysql"),
		host:       env.GetEnv("DATABASE_HOST", "localhost"),
		username:   env.GetEnv("DATABASE_USERNAME", "root"),
		password:   env.GetEnv("DATABASE_PASSWORD", ""),
		dbName:     env.GetEnv("DATABASE_NAME", ""),
		port:       env.GetEnv("DATABASE_PORT", "3306"),
	}
}
