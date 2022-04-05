package repository

import "github.com/go-pg/pg/v9/orm"

// QueryApply function signature for orm.Apply
type QueryApply func(query *orm.Query) (*orm.Query, error)
