package pager

import (
	"errors"
	"fmt"
	pkgerr "github.com/sanches1984/gopkg-pg-orm/errors"
	"math"

	"github.com/sanches1984/gopkg-pg-orm/repository"
	dbpager "github.com/sanches1984/gopkg-pg-orm/repository/pager"
)

var LastPageTag = pkgerr.NewTag()

// pager internal implementation of facade
type pager struct {
	Page       int32
	PageSize   int32
	TotalItems *int32
}

// GetOffset compute data offset
func (p *pager) GetOffset() int32 {
	return (p.Page - 1) * p.PageSize
}

// GetLimit compute data limit
func (p *pager) GetLimit() int32 {
	return p.PageSize
}

// GetPage get current page number
func (p *pager) GetPage() int32 {
	return p.Page
}

// GetPageSize get page size
func (p *pager) GetPageSize() int32 {
	return p.PageSize
}

// GetTotalPages compute total pages for all items
func (p *pager) GetTotalPages() int32 {
	if p.TotalItems == nil || *p.TotalItems < 1 {
		return 0
	}

	return int32(math.Ceil(float64(*p.TotalItems) / float64(p.PageSize)))
}

// GetTotalItems get total items
func (p *pager) GetTotalItems() int32 {
	if p.TotalItems == nil {
		return 0
	}
	return *p.TotalItems
}

// SetTotalItems set total items
func (p *pager) SetTotalItems(total int32) {
	p.TotalItems = &total
}

// NextPage try to step onto next page
func (p *pager) NextPage() error {
	if p.TotalItems == nil {
		return pkgerr.NewInternalError(errors.New("total items must be set"))
	}

	if p.Page*p.PageSize >= *p.TotalItems {
		return pkgerr.NewInternalError(fmt.Errorf("last page is %d", p.Page)).WithTag(LastPageTag)
	}

	p.Page++
	return nil
}

// GetApplyFn convert pager to query apply function
func (p *pager) GetApplyFn() repository.QueryApply {
	return dbpager.New(int(p.GetOffset()), int(p.GetLimit()))
}
