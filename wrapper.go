package database

import (
	"context"
	"io"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

type dbWrapper struct {
	ctx              context.Context
	Conn             *pg.DB
	Txn              *pg.Tx
	WrappedProcessor func(ctx context.Context, processor func() (orm.Result, error), query string, model interface{}) (orm.Result, error)
}

// NewDbClient ...
func NewDbClient(conn *pg.DB) IClient {
	return &dbWrapper{Conn: conn}
}

// Db ...
func (w *dbWrapper) Db() *pg.DB {
	return w.Conn
}

// Tx ...
func (w *dbWrapper) Tx() *pg.Tx {
	return w.Txn
}

// WrapWithContext ...
func (w *dbWrapper) WrapWithContext(ctx context.Context) IClient {
	return &dbWrapper{
		Conn:             w.Conn.WithContext(ctx),
		WrappedProcessor: w.WrappedProcessor,
	}
}

// StartTransaction ...
func (w *dbWrapper) StartTransaction() (IClient, error) {
	txn, err := w.Conn.Begin()
	if err != nil {
		return nil, err
	}

	return &dbWrapper{
		Conn:             w.Conn,
		Txn:              txn,
		WrappedProcessor: w.WrappedProcessor,
	}, nil
}

// SetWrappedQueryProcessor ...
func (w *dbWrapper) SetWrappedQueryProcessor(processor func(ctx context.Context, processor func() (orm.Result, error), query string, model interface{}) (orm.Result, error)) {
	w.WrappedProcessor = processor
}

// Context ...
func (w *dbWrapper) Context() context.Context {
	if w.Txn != nil {
		return w.Txn.Context()
	}
	return w.Conn.Context()
}

// WithContext ...
func (w *dbWrapper) WithContext(ctx context.Context) IClient {
	w.ctx = ctx
	return w
}

// Close ...
func (w *dbWrapper) Close() error {
	return w.Conn.Close()
}

// Model ...
func (w *dbWrapper) Model(model ...interface{}) *orm.Query {
	return orm.NewQuery(w, model...).Context(w.ctx)
}

// Select ...
func (w *dbWrapper) Select(model interface{}) error {
	if w.Txn != nil {
		return w.Txn.Select(model)
	}
	return w.Conn.Select(model)
}

// Insert ...
func (w *dbWrapper) Insert(model ...interface{}) error {
	if w.Txn != nil {
		return w.Txn.Insert(model...)
	}
	return w.Conn.Insert(model...)
}

// Update ...
func (w *dbWrapper) Update(model interface{}) error {
	if w.Txn != nil {
		return w.Txn.Update(model)
	}
	return w.Conn.Update(model)
}

// Delete ...
func (w *dbWrapper) Delete(model interface{}) error {
	if w.Txn != nil {
		return w.Txn.Delete(model)
	}
	return w.Conn.Delete(model)
}

// Exec ...
func (w *dbWrapper) Exec(query interface{}, params ...interface{}) (orm.Result, error) {
	processor := func() (orm.Result, error) {
		if w.Txn != nil {
			return w.Txn.Exec(query, params...)
		}
		return w.Conn.Exec(query, params...)
	}

	if w.WrappedProcessor == nil {
		return processor()
	}

	return w.WrappedProcessor(w.Conn.Context(), processor, w.queryString(query), nil)
}

// ExecOne ...
func (w *dbWrapper) ExecOne(query interface{}, params ...interface{}) (orm.Result, error) {
	res, err := w.Exec(query, params...)
	if err != nil {
		return nil, err
	}

	if err := w.assertOneRow(res.RowsAffected()); err != nil {
		return nil, err
	}
	return res, nil
}

// Query ...
func (w *dbWrapper) Query(model, query interface{}, params ...interface{}) (orm.Result, error) {
	processor := func() (orm.Result, error) {
		if w.Txn != nil {
			return w.Txn.Query(model, query, params...)
		}
		return w.Conn.Query(model, query, params...)
	}

	if w.WrappedProcessor == nil {
		return processor()
	}

	return w.WrappedProcessor(w.Conn.Context(), processor, w.queryString(query), model)
}

// QueryOne ...
func (w *dbWrapper) QueryOne(model, query interface{}, params ...interface{}) (orm.Result, error) {
	res, err := w.Query(model, query, params...)
	if err != nil {
		return nil, err
	}

	if err := w.assertOneRow(res.RowsAffected()); err != nil {
		return nil, err
	}
	return res, nil
}

// CopyFrom ...
func (w *dbWrapper) CopyFrom(r io.Reader, query interface{}, params ...interface{}) (orm.Result, error) {
	if w.Txn != nil {
		return w.Txn.CopyFrom(r, query, params...)
	}
	return w.Conn.CopyFrom(r, query, params...)
}

// CopyTo ...
func (w *dbWrapper) CopyTo(iw io.Writer, query interface{}, params ...interface{}) (orm.Result, error) {
	if w.Txn != nil {
		return w.Txn.CopyTo(iw, query, params...)
	}
	return w.Conn.CopyTo(iw, query, params...)
}

// FormatQuery ...
func (w *dbWrapper) FormatQuery(b []byte, query string, params ...interface{}) []byte {
	if w.Txn != nil {
		return w.Txn.Formatter().FormatQuery(b, query, params...)
	}
	return w.Conn.Formatter().FormatQuery(b, query, params...)
}

func (w *dbWrapper) assertOneRow(affected int) error {
	if affected == 0 {
		return pg.ErrNoRows
	}
	if affected > 1 {
		return pg.ErrMultiRows
	}
	return nil
}

func (w *dbWrapper) queryString(query interface{}) string {
	switch typed := query.(type) {
	case orm.QueryAppender:
		if b, err := typed.AppendQuery(w.Formatter(), nil); err == nil {
			return string(b)
		}
	case string:
		return typed
	}
	return ""
}

// ForceDelete ...
func (w *dbWrapper) ForceDelete(values interface{}) error {
	if w.Txn != nil {
		return w.Txn.ForceDelete(values)
	}
	return w.Conn.ForceDelete(values)
}

// ModelContext ...
func (w *dbWrapper) ModelContext(c context.Context, model ...interface{}) *orm.Query {
	if w.Txn != nil {
		return w.Txn.ModelContext(c, model...)
	}
	return w.Conn.ModelContext(c, model...)
}

// ExecContext ...
func (w *dbWrapper) ExecContext(c context.Context, query interface{}, params ...interface{}) (pg.Result, error) {
	if w.Txn != nil {
		return w.Txn.ExecContext(c, query, params...)
	}
	return w.Conn.ExecContext(c, query, params...)
}

// ExecOneContext ...
func (w *dbWrapper) ExecOneContext(c context.Context, query interface{}, params ...interface{}) (pg.Result, error) {
	if w.Txn != nil {
		return w.Txn.ExecOneContext(c, query, params...)
	}
	return w.Conn.ExecOneContext(c, query, params...)
}

// QueryContext ...
func (w *dbWrapper) QueryContext(c context.Context, model, query interface{}, params ...interface{}) (pg.Result, error) {
	if w.Txn != nil {
		return w.Txn.QueryContext(c, model, query, params...)
	}
	return w.Conn.QueryContext(c, model, query, params...)
}

// QueryOneContext ...
func (w *dbWrapper) QueryOneContext(c context.Context, model, query interface{}, params ...interface{}) (pg.Result, error) {
	if w.Txn != nil {
		return w.Txn.QueryOneContext(c, model, query, params...)
	}
	return w.Conn.QueryOneContext(c, model, query, params...)
}

// Formatter ...
func (w *dbWrapper) Formatter() orm.QueryFormatter {
	return w.Formatter()
}
