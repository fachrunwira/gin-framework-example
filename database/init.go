package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fachrunwira/go-query-builder/builder"
	_ "github.com/go-sql-driver/mysql"
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

	instance = db
	return nil
}

func Inject(ctx context.Context) context.Context {
	builder.SetContextKey(CtxValueDB)
	return context.WithValue(ctx, CtxValueDB, instance)
}

func Close() error {
	if instance != nil {
		return instance.Close()
	}
	return nil
}

func FromContext(ctx context.Context) (*sql.DB, error) {
	db, ok := ctx.Value(CtxValueDB).(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("database not found in context")
	}

	return db, nil
}

func GetDB() *sql.DB {
	return instance
}
