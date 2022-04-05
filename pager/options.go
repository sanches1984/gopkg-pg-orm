package pager

const (
	// DefaultPageSize default page size
	DefaultPageSize = 100
	// DefaultMaxPageSize default max page size
	DefaultMaxPageSize = 1000
)

// Options pager options
type Options struct {
	PageSize    int32
	MaxPageSize int32
}

// NewOptions create options with defaults
func NewOptions() *Options {
	return &Options{
		PageSize:    DefaultPageSize,
		MaxPageSize: DefaultMaxPageSize,
	}
}

// WithPageSize update options with new pageSize value
func (o *Options) WithPageSize(pageSize int32) *Options {
	o.PageSize = pageSize
	return o
}

// WithMaxPageSize update options with new maxPageSize value
func (o *Options) WithMaxPageSize(maxPageSize int32) *Options {
	o.MaxPageSize = maxPageSize
	return o
}
