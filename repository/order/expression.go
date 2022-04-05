package order

import "github.com/go-pg/pg/v9"

// Expression common facade for order expressions
type Expression interface {
	Expression() string
	Params() []interface{}
}

// Asc sort in ascending order
type Asc string

// AscNullsFirst sort in ascending order with NULL on top
type AscNullsFirst string

// AscNullsLast sort in ascending order with NULL at the end
type AscNullsLast string

// Desc sort in descending order
type Desc string

// DescNullsFirst sort in descending order with NULL on top
type DescNullsFirst string

// DescNullsLast sort in descending order with NULL at the end
type DescNullsLast string

// Expression provide query expression
func (e Asc) Expression() string {
	return "? ASC"
}

// Params provide query params
func (e Asc) Params() []interface{} {
	return []interface{}{
		pg.Ident(string(e)),
	}
}

// Expression provide query expression
func (e AscNullsFirst) Expression() string {
	return "? ASC NULLS FIRST"
}

// Params provide query params
func (e AscNullsFirst) Params() []interface{} {
	return []interface{}{
		pg.Ident(string(e)),
	}
}

// Expression provide query expression
func (e AscNullsLast) Expression() string {
	return "? ASC NULLS LAST"
}

// Params provide query params
func (e AscNullsLast) Params() []interface{} {
	return []interface{}{
		pg.Ident(string(e)),
	}
}

// Expression provide query expression
func (e Desc) Expression() string {
	return "? DESC"
}

// Params provide query params
func (e Desc) Params() []interface{} {
	return []interface{}{
		pg.Ident(string(e)),
	}
}

// Expression provide query expression
func (e DescNullsFirst) Expression() string {
	return "? DESC NULLS FIRST"
}

// Params provide query params
func (e DescNullsFirst) Params() []interface{} {
	return []interface{}{
		pg.Ident(string(e)),
	}
}

// Expression provide query expression
func (e DescNullsLast) Expression() string {
	return "? DESC NULLS LAST"
}

// Params provide query params
func (e DescNullsLast) Params() []interface{} {
	return []interface{}{
		pg.Ident(string(e)),
	}
}
