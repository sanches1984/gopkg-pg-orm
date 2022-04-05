package database

import (
	"context"
	"io"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

// IClient ...
type IClient interface {
	Db() *pg.DB
	Tx() *pg.Tx

	StartTransaction() (IClient, error)
	WrapWithContext(ctx context.Context) IClient
	SetWrappedQueryProcessor(func(ctx context.Context, processor func() (orm.Result, error), query string, model interface{}) (orm.Result, error))

	Context() context.Context
	WithContext(ctx context.Context) IClient
	Close() error

	Model(model ...interface{}) *orm.Query
	Select(model interface{}) error
	Insert(model ...interface{}) error
	Update(model interface{}) error
	Delete(model interface{}) error

	Exec(query interface{}, params ...interface{}) (orm.Result, error)
	ExecOne(query interface{}, params ...interface{}) (orm.Result, error)
	Query(model, query interface{}, params ...interface{}) (orm.Result, error)
	QueryOne(model, query interface{}, params ...interface{}) (orm.Result, error)

	CopyFrom(r io.Reader, query interface{}, params ...interface{}) (orm.Result, error)
	CopyTo(w io.Writer, query interface{}, params ...interface{}) (orm.Result, error)

	FormatQuery(b []byte, query string, params ...interface{}) []byte

	ForceDelete(model interface{}) error
}
