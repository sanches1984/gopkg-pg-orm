package database

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/rs/zerolog"
	"time"
)

const queryStartTime = "StartTime"

type IDBLogger interface {
	pg.QueryHook
}

type dbLogger struct {
	logger   zerolog.Logger
	duration time.Duration
}

func newDBLogger(logger zerolog.Logger, duration time.Duration) IDBLogger {
	return &dbLogger{
		logger:   logger,
		duration: duration,
	}
}

func (d *dbLogger) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	if event.Stash == nil {
		event.Stash = make(map[interface{}]interface{})
	}
	event.Stash[queryStartTime] = time.Now()
	return ctx, nil
}

func (d *dbLogger) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	query, err := event.FormattedQuery()
	if err == nil {
		var duration time.Duration
		if event.Stash != nil {
			if v, ok := event.Stash[queryStartTime]; ok {
				duration = time.Now().Sub(v.(time.Time))
			}
		}
		logEvent := d.logger.Info()
		if d.duration != 0 {
			if d.duration > duration {
				return nil
			}
			logEvent = d.logger.Warn()
		}
		txt := "query: " + query
		if duration != 0 {
			txt += fmt.Sprintf(" [%d ms]", duration.Nanoseconds()/1000000)
		}
		if event.Err != nil {
			txt += "\nerror: " + event.Err.Error()
		}

		logEvent.Msg(txt)
	}
	return nil
}
