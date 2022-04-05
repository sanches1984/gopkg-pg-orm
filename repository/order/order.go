package order

import "github.com/go-pg/pg/v9/orm"

// Order directions
const (
	DirAsc            = "ASC"
	DirAscNullsFirst  = "ASC NULLS FIRST"
	DirAscNullsLast   = "ASC NULLS LAST"
	DirDesc           = "DESC"
	DirDescNullsFirst = "DESC NULLS FIRST"
	DirDescNullsLast  = "DESC NULLS LAST"
)

// Order repository order, collection of expressions
type Order []Expression

// Apply update query with expressions from Order
func (s Order) Apply(query *orm.Query) (*orm.Query, error) {
	for _, expr := range s {
		query.OrderExpr(expr.Expression(), expr.Params()...)
	}
	return query, nil
}

// Expr construct Expression based on field name and order direction
func Expr(field, dir string) Expression {
	switch dir {
	case DirAsc:
		return Asc(field)
	case DirAscNullsFirst:
		return AscNullsFirst(field)
	case DirAscNullsLast:
		return AscNullsLast(field)
	case DirDesc:
		return Desc(field)
	case DirDescNullsFirst:
		return DescNullsFirst(field)
	case DirDescNullsLast:
		return DescNullsLast(field)
	}

	return nil
}
