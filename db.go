package database

import (
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"reflect"
)

var dbKey = "db"
var dbLoggerKey = "dbLogger"

// NewContext returns a new Context that carries db
func NewContext(ctx context.Context, client IClient, options ...Option) context.Context {
	ctx = context.WithValue(ctx, &dbKey, client)
	for _, f := range options {
		ctx = f(ctx)
	}
	return ctx
}

// FromContext returns the DB value stored in ctx
func FromContext(ctx context.Context) IClient {
	client, ok := ctx.Value(&dbKey).(IClient)
	if !ok {
		return nil
	}
	return client.WithContext(ctx)
}

// Connect ...
func Connect(AppName string, cfg *pg.Options) IClient {
	if cfg.OnConnect == nil {
		cfg.OnConnect = onConnect(AppName)
	}

	return NewDbClient(pg.Connect(cfg))
}

// set client timezone to UTC
func onConnect(appName string) func(conn *pg.Conn) error {
	return func(conn *pg.Conn) error {
		conn.Exec("set timezone='UTC'")
		conn.Exec("set application_name=?", appName)
		return nil
	}
}

// GetTableName returns table name by model
func GetTableName(model interface{}) string {
	if t := orm.GetTable(reflect.TypeOf(model)); t != nil {
		return t.Name
	}
	return ""
}

type transactFunc func(client IClient) error

// PerformTransaction ...
func PerformTransaction(client IClient, fns ...transactFunc) error {
	txClient, err := client.StartTransaction()
	if err != nil {
		return err
	}

	tx := txClient.Tx()
	for _, fn := range fns {
		if err := fn(txClient); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
