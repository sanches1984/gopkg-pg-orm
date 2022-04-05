package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	pkgerr "github.com/sanches1984/gopkg-pg-orm/errors"
	"reflect"
	"strings"
	"time"

	db "github.com/sanches1984/gopkg-pg-orm"
	"github.com/sanches1984/gopkg-pg-orm/repository/opt"

	"github.com/go-pg/pg/v9/orm"
)

// DAO is a data access object
type DAO struct{}

// New creates new DAO structure
func New() *DAO {
	return &DAO{}
}

// DeletedSetter is an interface
type DeletedSetter interface {
	SetDeleted(time.Time)
}

func (r *DAO) Ping(ctx context.Context) error {
	_, err := db.FromContext(ctx).Exec("SELECT 1")
	return err
}

// WithTX executes passed function within transaction
func (r *DAO) WithTX(ctx context.Context, fn func(context.Context) error) error {
	dbc := db.FromContext(ctx)
	if dbc.Tx() != nil {
		return fn(ctx)
	}

	dbcTx, err := dbc.StartTransaction()
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}

	ctxTx := db.NewContext(ctx, dbcTx)
	err = fn(ctxTx)
	if err != nil {
		rollbackErr := dbcTx.Tx().Rollback()
		if rollbackErr != nil {
			// todo get logger from context
			log.Error().Err(pkgerr.Convert(ctx, rollbackErr)).Msg("failed to rollback transaction")
		}

		return err
	}

	err = dbcTx.Tx().Commit()
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}
	return nil
}

// FindOne selects the only record from database according to opts
func (r *DAO) FindOne(ctx context.Context, receiver interface{}, opts []opt.FnOpt) error {
	err := db.FromContext(ctx).Model(receiver).Apply(opt.Apply(opts...)).First()
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}

	return nil
}

// FindList selects all records from database according to opts
func (r *DAO) FindList(ctx context.Context, receiver interface{}, opts []opt.FnOpt) error {
	q := db.FromContext(ctx).Model(receiver).Apply(opt.ApplyFilter(opts...))

	err := q.Apply(opt.ApplyPaging(opts...)).Select()
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}

	return nil
}

// FindListWithTotal selects all records and total count of records from database according to opts
func (r *DAO) FindListWithTotal(ctx context.Context, receiver interface{}, opts []opt.FnOpt) (int, error) {
	q := db.FromContext(ctx).Model(receiver).Apply(opt.ApplyFilter(opts...))
	total, err := q.Count()
	if err != nil {
		return 0, pkgerr.Convert(ctx, err)
	}

	err = q.Apply(opt.ApplyPaging(opts...)).Select()
	if err != nil {
		return 0, pkgerr.Convert(ctx, err)
	}

	return total, nil
}

// GetTotal get total count of records from database according to opts
func (r *DAO) GetTotal(ctx context.Context, receiver interface{}, opts []opt.FnOpt) (int, error) {
	q := db.FromContext(ctx).Model(receiver).Apply(opt.ApplyFilter(opts...))
	total, err := q.Count()
	if err != nil {
		return 0, pkgerr.Convert(ctx, err)
	}

	return total, nil
}

// Update updates a record
func (r *DAO) Update(ctx context.Context, rec interface{}, columns ...string) error {
	columns = append(columns, "updated")
	q := db.FromContext(ctx).Model(rec).Column(columns...)
	// Slice not require additional filter
	if reflect.ValueOf(rec).Elem().Type().Kind() != reflect.Slice {
		q.WherePK()
	}
	_, err := q.Update()
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}

	return nil
}

// UpdateWhere updates a record with condition
func (r *DAO) UpdateWhere(ctx context.Context, rec interface{}, opts []opt.FnOpt, setFieldValuePairs ...interface{}) error {
	if len(setFieldValuePairs)&1 != 0 {
		return pkgerr.NewInternalError(fmt.Errorf("UpdateWhere: setFieldValuePairs must be even, got %d", len(setFieldValuePairs)))
	}
	setFieldValuePairs = append(setFieldValuePairs, "updated", time.Now())
	q := db.FromContext(ctx).Model(rec).Apply(opt.ApplyFilter(opts...))
	for i := 0; i < len(setFieldValuePairs); i += 2 {
		column, ok := setFieldValuePairs[i].(string)
		if !ok {
			return pkgerr.NewInternalError(fmt.Errorf("UpdateWhere: field must be string, got %T (%v)", setFieldValuePairs[i], setFieldValuePairs[i]))
		}
		q.Set(column+" = ?", setFieldValuePairs[i+1])
	}
	_, err := q.Update()
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}

	return nil
}

// UpdateWithReturning updates a record
func (r *DAO) UpdateWithReturning(ctx context.Context, rec interface{}, columns ...string) error {
	columns = append(columns, "updated")
	_, err := db.FromContext(ctx).Model(rec).Column(columns...).WherePK().Returning("*").Update()
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}

	return nil
}

// Insert creates a new record
func (r *DAO) Insert(ctx context.Context, rec ...interface{}) error {
	err := db.FromContext(ctx).Insert(rec...)
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}

	return nil
}

// SoftDelete marks record as deleted
func (r *DAO) SoftDelete(ctx context.Context, rec DeletedSetter) error {
	rec.SetDeleted(time.Now())
	err := r.Update(ctx, rec, "deleted")
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}

	return nil
}

// HardDelete removes record from database
func (r *DAO) HardDelete(ctx context.Context, rec interface{}) error {
	err := db.FromContext(ctx).Delete(rec)
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}

	return nil
}

// HardDeleteWhere removes record from database
func (r *DAO) HardDeleteWhere(ctx context.Context, rec interface{}, opts []opt.FnOpt) error {
	_, err := db.FromContext(ctx).Model(rec).Apply(opt.ApplyFilter(opts...)).Delete()
	if err != nil {
		return pkgerr.Convert(ctx, err)
	}

	return nil
}

// Upsert inserts recs, on conflict update columns
func (r *DAO) Upsert(ctx context.Context, recs interface{}, keys []string, columns ...string) error {
	if len(keys) == 0 {
		return pkgerr.NewBadRequestError(errors.New("keys cannot be empty"))
	}

	goNames := make([]string, 0, len(keys))
	if t := orm.GetTable(getType(recs)); t != nil {
		for _, key := range keys {
			goNames = append(goNames, t.FieldsMap[key].GoName)
		}
	}

	var models []interface{}
	k := reflect.TypeOf(recs).Kind()
	if k == reflect.Slice {
		models = GetUniqueModels(recs, func(model interface{}) string {
			values := []string{}
			for _, key := range goNames {
				values = append(values, fmt.Sprint(reflect.ValueOf(model).Elem().FieldByName(key)))
			}
			return strings.Join(values, "_")
		})
	} else if k == reflect.Ptr {
		models = []interface{}{recs}
	} else {
		return pkgerr.NewBadRequestError(errors.New("recs must be slice or pointer to struct"))
	}

	if len(models) == 0 {
		return pkgerr.NewBadRequestError(errors.New("models cannot be empty"))
	}

	dbc := db.FromContext(ctx)
	q := dbc.Model(&models).OnConflict("(" + strings.Join(keys, ",") + ") DO UPDATE")

	for _, column := range columns {
		q = q.Set(column + " = EXCLUDED." + column)
	}

	_, err := q.Insert()
	return err
}

func getType(models interface{}) reflect.Type {
	var m interface{}

	switch reflect.TypeOf(models).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(models)
		if s.Index(0).Kind() == reflect.Ptr {
			m = s.Index(0).Elem().Interface()
		} else {
			m = s.Index(0).Interface()
		}

	case reflect.Ptr:
		m = reflect.ValueOf(models).Elem().Interface()

	case reflect.Struct:
		m = reflect.ValueOf(models).Interface()

	}

	return reflect.TypeOf(m)
}

// GetUniqueModels - make models unique according to key returned by f
// if two models have the same key, the last one takes precedence
func GetUniqueModels(models interface{}, f func(model interface{}) string) []interface{} {
	rows := make(map[string]interface{})

	switch reflect.TypeOf(models).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(models)

		for i := 0; i < s.Len(); i++ {
			var model interface{}
			if s.Index(i).Kind() == reflect.Ptr {
				model = s.Index(i).Interface()
			} else {
				model = s.Index(i).Addr().Interface()
			}

			rows[f(model)] = model
		}
	}

	unique := make([]interface{}, 0, len(rows))
	for i := range rows {
		unique = append(unique, rows[i])
	}

	return unique
}
