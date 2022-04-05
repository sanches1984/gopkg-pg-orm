package opt

import (
	"reflect"
	"strings"

	"github.com/sanches1984/gopkg-pg-orm/pager"
	"github.com/sanches1984/gopkg-pg-orm/repository"
	"github.com/sanches1984/gopkg-pg-orm/repository/filter"
	"github.com/sanches1984/gopkg-pg-orm/repository/order"

	"github.com/go-pg/pg/v9/orm"
)

// Opt is options for database requests
type Opt struct {
	Page      int32
	PageSize  int32
	SortBy    string
	SortOrder string
	Filter    filter.Filter
	Fn        []repository.QueryApply
}

// FnOpt is a function that modifies options
type FnOpt func(*Opt)

// New creates new Opt
func New(optFn ...FnOpt) *Opt {
	o := &Opt{}
	for _, fn := range optFn {
		if fn != nil {
			fn(o)
		}
	}

	return o
}

// Apply returns a function that builds request with all statements
func (o *Opt) Apply() repository.QueryApply {
	return func(query *orm.Query) (*orm.Query, error) {
		if o == nil {
			return query, nil
		}

		query, _ = o.ApplyFilter()(query)
		query, _ = o.ApplyFn()(query)
		query, _ = o.ApplyPaging()(query)
		return query, nil
	}
}

// ApplyFilter returns a function that builds request only with WHERE statements
func (o *Opt) ApplyFilter() repository.QueryApply {
	return func(query *orm.Query) (*orm.Query, error) {
		if o == nil {
			return query, nil
		}

		if o.IsFilter() {
			query, _ = o.Filter.Apply(query)
		}

		return query, nil
	}
}

// ApplyPaging returns a function that builds request only with ORDER BY, LIMIT ... OFFSET statements
func (o *Opt) ApplyPaging() repository.QueryApply {
	return func(query *orm.Query) (*orm.Query, error) {
		if o == nil {
			return query, nil
		}

		if o.IsPaging() {
			query = query.Apply(pager.NewPagerWithPageSize(o.Page, o.PageSize).GetApplyFn())
		}

		if o.IsSorting() {
			query = query.Apply(order.Order{order.Expr(o.SortBy, o.SortOrder)}.Apply)
		}

		return query, nil
	}
}

// ApplyFn returns a function that builds request only with WHERE statements
func (o *Opt) ApplyFn() repository.QueryApply {
	return func(query *orm.Query) (*orm.Query, error) {
		if o == nil {
			return query, nil
		}

		if o.IsFn() {
			for _, fn := range o.Fn {
				query = query.Apply(fn)
			}
		}

		return query, nil
	}
}

// ApplyFilter calls ApplyFilter for each FnOpt in a chain
func ApplyFilter(optFn ...FnOpt) repository.QueryApply {
	return New(optFn...).ApplyFilter()
}

// ApplyPaging calls ApplyPaging for each FnOpt in a chain
func ApplyPaging(optFn ...FnOpt) repository.QueryApply {
	return New(optFn...).ApplyPaging()
}

// Apply calls Apply for each FnOpt in a chain
func Apply(optFn ...FnOpt) repository.QueryApply {
	return New(optFn...).Apply()
}

// List converts periodic opts args into slice
func List(optFn ...FnOpt) []FnOpt {
	return optFn
}

// IsFn responds whether fn options set
func (o *Opt) IsFn() bool {
	return len(o.Fn) > 0
}

// IsPaging responds whether pagination options set
func (o *Opt) IsPaging() bool {
	return o.PageSize > 0 || o.Page > 0
}

// IsSorting responds whether sorting options set
func (o *Opt) IsSorting() bool {
	return o.SortOrder != "" && o.SortBy != ""
}

// IsFilter responds whether filter options set
func (o *Opt) IsFilter() bool {
	return len(o.Filter) > 0
}

// Add query function
func Fn(queryFn ...repository.QueryApply) FnOpt {
	return func(opt *Opt) {
		opt.Fn = append(opt.Fn, queryFn...)
	}
}

// Page sets page option
func Page(page int32) FnOpt {
	return func(opt *Opt) {
		opt.Page = page
	}
}

// PageSize sets page size option
func PageSize(size int32) FnOpt {
	return func(opt *Opt) {
		opt.PageSize = size
	}
}

// Limit sets page size option
func Limit(size int32) FnOpt {
	return PageSize(size)
}

// Paging sets both page and page size options
func Paging(page, size int32) FnOpt {
	return func(opt *Opt) {
		opt.Page = page
		opt.PageSize = size
	}
}

// Asc sets ascending order options
func Asc(column string) FnOpt {
	return func(opt *Opt) {
		opt.SortBy = column
		opt.SortOrder = order.DirAsc
	}
}

// Desc sets descending order options
func Desc(column string) FnOpt {
	return func(opt *Opt) {
		opt.SortBy = column
		opt.SortOrder = order.DirDesc
	}
}

// Order options
func Order(columnAndDirection string) FnOpt {
	return func(opt *Opt) {
		arr := strings.SplitN(columnAndDirection, " ", 2)
		opt.SortBy = arr[0]
		if len(arr) == 1 {
			opt.SortOrder = order.DirAsc
		} else {
			switch strings.ToUpper(strings.TrimSpace(arr[1])) {
			case "ASC":
				opt.SortOrder = order.DirAsc
			case "ASC NULLS FIRST":
				opt.SortOrder = order.DirAscNullsFirst
			case "ASC NULLS LAST":
				opt.SortOrder = order.DirAscNullsLast
			case "DESC":
				opt.SortOrder = order.DirDesc
			case "DESC NULLS FIRST":
				opt.SortOrder = order.DirDescNullsFirst
			case "DESC NULLS LAST":
				opt.SortOrder = order.DirDescNullsLast
			default:
				panic("Unknown order rule: " + arr[1])
			}
		}
	}
}

// Eq adds to filter equal condition
func Eq(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Eq{column: val})
	}
}

// Eq adds to filter equal in lower case condition
func EqLower(column, val string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.EqLower{column: val})
	}
}

// Gt adds to filter great than condition
func Gt(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Gt{column: val})
	}
}

// Ge adds to filter great and equal condition
func Ge(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Ge{column: val})
	}
}

// Lt adds to filter less than condition
func Lt(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Lt{column: val})
	}
}

// Le adds to filter less and equal condition
func Le(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Le{column: val})
	}
}

// Between adds to filter between condition
func Between(column string, from, to interface{}) FnOpt {
	return func(opt *Opt) {
		vals := make([]interface{}, 2)
		vals[0] = from
		vals[1] = to
		opt.Filter = append(opt.Filter, filter.Between{column: vals})
	}
}

// Neq adds to filter not-equal condition
func Neq(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Ne{column: val})
	}
}

// MayIn sets condition for IN operation only if vals is not empty
func MayIn(column string, vals interface{}) FnOpt {
	if reflect.TypeOf(vals).Kind() == reflect.Slice && reflect.ValueOf(vals).Len() == 0 {
		return nil
	}

	return In(column, vals)
}

// In sets condition for IN operation
func In(column string, vals interface{}) FnOpt {
	return func(opt *Opt) {
		if reflect.TypeOf(vals).Kind() != reflect.Slice {
			vals = []interface{}{vals}
		}

		in := []interface{}{}
		v := reflect.ValueOf(vals)
		for i := 0; i < v.Len(); i++ {
			in = append(in, v.Index(i).Interface())
		}

		opt.Filter = append(opt.Filter, filter.In{column: in})
	}
}

// Contains builds a condition with `LIKE %val%` statement
func Contains(column string, val string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Contains{column: val})
	}
}

// Starts builds a condition with `LIKE val%` statement
func Starts(column string, val string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Starts{column: val})
	}
}

// Ends builds a condition with `LIKE %val` statement
func Ends(column string, val string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Ends{column: val})
	}
}

// JsonEq builds a condition with `path.to.value = val` statement
func JsonEq(column string, val interface{}, path ...string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.JsonEq{
			Column: column,
			Path:   path,
			Value:  val,
		})
	}
}

// JsonContains builds a condition with `path.to.value @> val` statement
func JsonContains(column string, val ...filter.JsonPath) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.JsonContains{
			Column: column,
			Path:   val,
		})
	}
}

// JsonContainsValue builds a condition with `path.to.value LIKE %val%` statement
func JsonContainsValue(column string, val string, path ...string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.JsonContainsValue{
			Column: column,
			Path:   path,
			Value:  val,
		})
	}
}

// Or adds set of conditions joined with OR statement
func Or(optFn ...FnOpt) FnOpt {
	return func(opt *Opt) {
		o := New(optFn...)
		opt.Filter = append(opt.Filter, filter.Or(o.Filter))
	}
}

// And adds set of conditions joined with AND statement
func And(optFn ...FnOpt) FnOpt {
	return func(opt *Opt) {
		o := New(optFn...)
		opt.Filter = append(opt.Filter, filter.And(o.Filter))
	}
}

// NotNull adds `IS NOT NULL` condition
func NotNull(column string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.NotNull(column))
	}
}

// IsNull adds `IS NULL` condition
func IsNull(column string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.IsNull(column))
	}
}

// Not adds `NOT` condition
func Not(optFn ...FnOpt) FnOpt {
	return func(opt *Opt) {
		o := New(optFn...)
		opt.Filter = append(opt.Filter, filter.Not(o.Filter))
	}
}
