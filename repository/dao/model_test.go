//+build !ci

package dao

import (
	"context"
	"time"
)

// Agent is a test model
type Agent struct {
	tableName    struct{}   `pg:"agent"`
	ID           int64      `pg:"id,unique"`
	Name         string     `pg:"name,notnull,use_zero"`
	State        string     `pg:"state,notnull,use_zero"`
	ServiceLevel *string    `pg:"service_level"`
	INN          string     `pg:"inn"`
	Meta         string     `pg:"meta"`
	IsBlocked    bool       `pg:"is_blocked,notnull,use_zero"`
	Created      time.Time  `pg:"created,notnull,type:timestamp,default:now()"`
	Updated      time.Time  `pg:"updated,notnull,type:timestamp,default:now()"`
	Deleted      *time.Time `pg:"deleted,type:timestamp"`
}

const (
	AgentStateRegistered string = "registered"
	AgentStateApproved   string = "approved"
)

// SetDeletedAt sets deleted field
func (b *Agent) SetDeleted(t time.Time) {
	b.Deleted = &t
}

// BeforeInsert is a callback
func (b *Agent) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	b.Created = now
	b.Updated = now
	if b.State == "" {
		b.State = AgentStateRegistered
	}
	return ctx, nil
}

// BeforeUpdate is a callback
func (b *Agent) BeforeUpdate(ctx context.Context) (context.Context, error) {
	b.Updated = time.Now()
	return ctx, nil
}
