package pager

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/sanches1984/gopkg-database/repository"
)

// pager implementation for orm.Apply
type pager struct {
	offset int
	limit  int
}

// New construct pager from raw values
func New(offset, limit int) repository.QueryApply {
	return pager{offset, limit}.apply
}

// apply implementation of repository.ApplyFn
func (p pager) apply(query *orm.Query) (*orm.Query, error) {
	return query.Offset(p.offset).Limit(p.limit), nil
}
