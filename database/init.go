package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DatabaseOptions struct {
	MaxOpenConnection     int
	MaxIdleConnection     int
	MaxConnectionLifetime time.Duration
}

type dbKey string

const CtxValueDB dbKey = "database"

var instance *sql.DB

func Init(options *DatabaseOptions) error {
	conf := getConfig()

	dsn := mysql.NewConfig()
	dsn.User = conf.username
	dsn.DBName = conf.dbName
	dsn.Addr = fmt.Sprintf("%s:%s", conf.host, conf.port)
	dsn.Passwd = conf.password
	dsn.ParseTime = true

	db, err := sql.Open(conf.connection, dsn.FormatDSN())
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	if options != nil {
		db.SetMaxOpenConns(options.MaxOpenConnection)
		db.SetMaxIdleConns(options.MaxIdleConnection)
		db.SetConnMaxLifetime(options.MaxConnectionLifetime)
	} else {
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(2 * time.Minute)
	}

	instance = db
	return nil
}

func Close() error {
	if instance != nil {
		return instance.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return instance
}
