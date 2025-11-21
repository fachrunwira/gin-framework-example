package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseOptions struct {
	MaxOpenConnection     int
	MaxIdleConnection     int
	MaxConnectionLifetime time.Duration
}

type databaseCon struct {
	*sql.DB
}

type dbKey string

const ctxDbValue dbKey = "database"

var instance *databaseCon

func Init(options *DatabaseOptions) error {
	conf := getConfig()

	db, err := sql.Open(conf.connection, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.username, conf.password, conf.host, conf.port, conf.dbName))
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

	instance = &databaseCon{db}
	return nil
}

func Inject(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxDbValue, instance)
}

func Close() error {
	if instance != nil {
		return instance.DB.Close()
	}
	return nil
}

func FromContext(ctx context.Context) (*sql.DB, error) {
	db, ok := ctx.Value(ctxDbValue).(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("database not found in context")
	}

	return db, nil
}

func GetInstance() *databaseCon {
	return instance
}

func GetDB() *sql.DB {
	return instance.DB
}
