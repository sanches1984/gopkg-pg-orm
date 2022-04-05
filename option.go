package database

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

type Option func(ctx context.Context) context.Context

func WithLogger(logger zerolog.Logger, duration time.Duration) Option {
	logger.Info().Msgf("db logger: log query over %v", duration)
	return func(ctx context.Context) context.Context {
		dbLogger := newDBLogger(logger, duration)
		dbc := FromContext(ctx)
		if dbc == nil {
			return ctx
		}
		dbc.Db().AddQueryHook(dbLogger)
		return context.WithValue(ctx, &dbLoggerKey, dbLogger)
	}
}
