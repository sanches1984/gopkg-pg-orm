package filter

import "github.com/go-pg/pg/v9/orm"

// Filter repository filter, collection of conditions
type Filter []Condition

// Apply update query with conditions from filter
func (f Filter) Apply(query *orm.Query) (*orm.Query, error) {
	for _, cond := range f {
		query.Where(cond.Condition(), cond.Params()...)
	}
	return query, nil
}
