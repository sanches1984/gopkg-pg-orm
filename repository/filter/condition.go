package filter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-pg/pg/v9"
)

// JsonPath store json path and value
type JsonPath struct {
	Path  string
	Value interface{}
}

// Condition common facade for filter conditions
type Condition interface {
	Condition() string
	Params() []interface{}
}

// Eq field name equal to value
type Eq map[string]interface{}

// Eq field name equal to value
type EqLower map[string]string

// Between field name equal to value
type Between map[string]interface{}

// Ne field name not equal to value
type Ne map[string]interface{}

// Lt field name less than value
type Lt map[string]interface{}

// Le field name less than value or equal
type Le map[string]interface{}

// Gt field name greater than value
type Gt map[string]interface{}

// Ge field name greater than value or equal
type Ge map[string]interface{}

// In field name contains in list of values
type In map[string][]interface{}

// InInt64 field name contains in list of ints
type InInt64 map[string][]int64

// InStr field name contains in list of strings
type InStr map[string][]string

// NotIn field name not contains in list of values
type NotIn map[string][]interface{}

// NotInStr field name not contains in list of strings
type NotInStr map[string][]string

// Starts filter
type Starts map[string]string

// Contains filter
type Contains map[string]string

// Ends filter
type Ends map[string]string

// Match filter
type Match map[string]string

// MatchMany filter
type MatchMany map[string][]string

// IsNull field name equal to NULL
type IsNull string

// NotNull field name not equal to NULL
type NotNull string

// Or filter
type Or []Condition

// Or filter
type And []Condition

// Not filter
type Not []Condition

// Raw ...
type Raw struct {
	Query       string
	QueryParams []interface{}
}

// JsonEq filter
type JsonEq struct {
	Column string
	Path   []string
	Value  interface{}
}

// JsonContains filter
type JsonContains struct {
	Column string
	Path   []JsonPath
}

// JsonContainsValue filter
type JsonContainsValue struct {
	Column string
	Path   []string
	Value  interface{}
}

// Condition provide query condition
func (c Eq) Condition() string {
	return "? = ?"
}

// Params provide query params
func (c Eq) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			val,
		}
	}
	return nil
}

// Condition provide query condition
func (c EqLower) Condition() string {
	return "LOWER(?) = ?"
}

// Params provide query params
func (c EqLower) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			strings.ToLower(val),
		}
	}
	return nil
}

// Condition provide query condition
func (c Between) Condition() string {
	return "? BETWEEN ? AND ?"
}

// Params provide query params
func (c Between) Params() []interface{} {
	for key, vals := range c {
		v := vals.([]interface{})
		ret := []interface{}{}
		ret = append(ret, pg.Ident(key))
		return append(ret, v...)
	}
	return nil
}

// Condition provide query condition
func (c Ne) Condition() string {
	return "? != ?"
}

// Params provide query params
func (c Ne) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			val,
		}
	}
	return nil
}

// Condition provide query condition
func (c Lt) Condition() string {
	return "? < ?"
}

// Params provide query params
func (c Lt) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			val,
		}
	}
	return nil
}

// Condition provide query condition
func (c Le) Condition() string {
	return "? <= ?"
}

// Params provide query params
func (c Le) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			val,
		}
	}
	return nil
}

// Condition provide query condition
func (c Gt) Condition() string {
	return "? > ?"
}

// Params provide query params
func (c Gt) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			val,
		}
	}
	return nil
}

// Condition provide query condition
func (c Ge) Condition() string {
	return "? >= ?"
}

// Params provide query params
func (c Ge) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			val,
		}
	}
	return nil
}

// Condition provide query condition
func (c In) Condition() string {
	return "? IN (?)"
}

// Params provide query params
func (c In) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			pg.In(val),
		}
	}
	return nil
}

// Condition provide query condition
func (c InInt64) Condition() string {
	return "? IN (?)"
}

// Params provide query params
func (c InInt64) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			pg.In(val),
		}
	}
	return nil
}

// Condition provide query condition
func (c InStr) Condition() string {
	return "? IN (?)"
}

// Params provide query params
func (c InStr) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			pg.In(val),
		}
	}
	return nil
}

// Condition provide query condition
func (c NotIn) Condition() string {
	return "? NOT IN (?)"
}

// Params provide query params
func (c NotIn) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			pg.In(val),
		}
	}
	return nil
}

// Condition provide query condition
func (c NotInStr) Condition() string {
	return "? NOT IN (?)"
}

// Params provide query params
func (c NotInStr) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			pg.In(val),
		}
	}
	return nil
}

// Condition provide query condition
func (c Starts) Condition() string {
	return "? ILIKE ?"
}

// Params provide query params
func (c Starts) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			val + "%",
		}
	}
	return nil
}

// Condition provide query condition
func (c Contains) Condition() string {
	return "CAST(? AS text) ILIKE ?"
}

// Params provide query params
func (c Contains) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			"%" + val + "%",
		}
	}
	return nil
}

// Condition provide query condition
func (c Ends) Condition() string {
	return "? ILIKE ?"
}

// Params provide query params
func (c Ends) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			"%" + val,
		}
	}
	return nil
}

// Condition provide query condition
func (c Match) Condition() string {
	return "to_tsvector('russian',?) @@ plainto_tsquery('russian',?)"
}

// Params provide query params
func (c Match) Params() []interface{} {
	for key, val := range c {
		return []interface{}{
			pg.Ident(key),
			val,
		}
	}
	return nil
}

// Condition provide query condition
func (c MatchMany) Condition() string {
	cnt := 1
	for _, val := range c {
		cnt = len(val)
		break
	}

	return "plainto_tsquery('russian',?) @@ (to_tsvector('russian',?)" + strings.Repeat("||to_tsvector('russian',?)", cnt-1) + ")"
}

// Params provide query params
func (c MatchMany) Params() []interface{} {
	for val, keys := range c {
		result := []interface{}{val}
		for _, key := range keys {
			result = append(result, pg.Ident(key))
		}
		return result
	}
	return nil
}

// Condition provide query condition
func (c IsNull) Condition() string {
	return "? IS NULL"
}

// Params provide query params
func (c IsNull) Params() []interface{} {
	return []interface{}{
		pg.Ident(string(c)),
	}
}

// Condition provide query condition
func (c NotNull) Condition() string {
	return "? IS NOT NULL"
}

// Params provide query params
func (c NotNull) Params() []interface{} {
	return []interface{}{
		pg.Ident(string(c)),
	}
}

// Condition provide query condition
func (c JsonContains) Condition() string {
	return "? @> ?"
}

// Params provide query params
func (c JsonContains) Params() []interface{} {
	check := make(map[string]interface{}, len(c.Path))
	for _, item := range c.Path {
		check[item.Path] = item.Value
	}
	str, _ := json.Marshal(check)
	return []interface{}{
		pg.Ident(c.Column),
		string(str),
	}
}

// Condition provide query condition
func (c JsonEq) Condition() string {
	if reflect.TypeOf(c.Value).Kind() == reflect.Slice {
		switch reflect.TypeOf(c.Value).Elem().Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16:
			return "(? #>> ?)::int IN (?)"
		case reflect.Int64, reflect.Uint64:
			return "(? #>> ?)::bigint IN (?)"
		case reflect.Float32, reflect.Float64:
			return "(? #>> ?)::float IN (?)"
		default:
			return "? #>> ? IN (?)"
		}
	}
	return "? #>> ? = ?"
}

// Params provide query params
func (c JsonEq) Params() []interface{} {
	var value interface{}
	kind := reflect.TypeOf(c.Value).Kind()
	switch kind {
	case reflect.String,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Int64, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		value = fmt.Sprintf("%v", c.Value)
	case reflect.Slice:
		value = pg.In(c.Value)
		kind2 := reflect.TypeOf(c.Value).Elem().Kind()
		switch kind2 {
		case reflect.Int64:

		default:
			panic("Unknown JsonEq slice type: " + kind2.String())
		}
	default:
		panic("Unknown JsonEq type: " + kind.String())
	}
	return []interface{}{
		pg.Ident(c.Column),
		"{" + strings.Join(c.Path, ",") + "}",
		value,
	}
}

// Condition provide query condition
func (c JsonContainsValue) Condition() string {
	return "? #>> ? ILIKE ?"
}

// Params provide query params
func (c JsonContainsValue) Params() []interface{} {
	return []interface{}{
		pg.Ident(c.Column),
		"{" + strings.Join(c.Path, ",") + "}",
		fmt.Sprintf("%%%v%%", c.Value),
	}
}

// Condition provide query condition
func (c Raw) Condition() string {
	return c.Query
}

// Params provide query params
func (c Raw) Params() []interface{} {
	return c.QueryParams
}

// Condition provide query condition
func (c Or) Condition() string {
	query := "("
	for i, cond := range c {
		if i > 0 {
			query += " OR "
		}
		query += fmt.Sprintf("(%v)", cond.Condition())
	}

	return query + ")"
}

// Params provide query params
func (c Or) Params() []interface{} {
	var result []interface{}
	for _, cond := range c {
		result = append(result, cond.Params()...)
	}
	return result
}

// Condition provide query condition
func (c And) Condition() string {
	query := "("
	for i, cond := range c {
		if i > 0 {
			query += " AND "
		}
		query += fmt.Sprintf("(%v)", cond.Condition())
	}

	return query + ")"
}

// Params provide query params
func (c And) Params() []interface{} {
	var result []interface{}
	for _, cond := range c {
		result = append(result, cond.Params()...)
	}
	return result
}

// Condition provide query condition
func (c Not) Condition() string {
	query := "NOT ("
	for i, cond := range c {
		if i > 0 {
			query += " AND "
		}
		query += fmt.Sprintf("(%v)", cond.Condition())
	}

	return query + ")"
}

// Params provide query params
func (c Not) Params() []interface{} {
	var result []interface{}
	for _, cond := range c {
		result = append(result, cond.Params()...)
	}
	return result
}
