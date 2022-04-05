package pager

import "github.com/sanches1984/gopkg-pg-orm/repository"

// Pager facade for pager
type Pager interface {
	GetOffset() int32
	GetLimit() int32
	GetPage() int32
	GetPageSize() int32
	GetTotalPages() int32
	GetTotalItems() int32
	SetTotalItems(total int32)
	NextPage() error
	GetApplyFn() repository.QueryApply
}

// NewPager construct new pager according to options
func NewPager(page int32, opt *Options) Pager {
	if page < 1 {
		page = 1
	}

	pageSize := opt.PageSize
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > opt.MaxPageSize {
		pageSize = opt.MaxPageSize
	}

	return &pager{
		Page:     page,
		PageSize: pageSize,
	}
}

// NewPagerWithPageSize construct new pager with pageSize option
func NewPagerWithPageSize(page int32, pageSize int32) Pager {
	return NewPager(page, NewOptions().WithPageSize(pageSize))
}

type pageGetter interface {
	GetPage() int32
	GetPageSize() int32
}

// NewRequestPager construct new pager from request
func NewRequestPager(req pageGetter) Pager {
	return NewPager(req.GetPage(), NewOptions().WithPageSize(req.GetPageSize()))
}
