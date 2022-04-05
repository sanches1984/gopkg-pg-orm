package errors

import (
	"context"
	"errors"
	"strings"

	"github.com/go-pg/pg/v9"
)

const (
	pgDuplicateErr = "duplicate key value"
	pgCodeField    = 'C'
	pgStatusField  = 'S'
	pgMessageField = 'M'
)

// Convert ...
func Convert(ctx context.Context, err error) Error {
	for {
		if err == pg.ErrNoRows {
			return NewNotFoundError(err)
		} else if err == pg.ErrMultiRows {
			return NewBadRequestError(err)
		}

		if errTyped, ok := err.(pg.Error); ok {
			return convert(errTyped)
		}

		oldErr := err
		err = errors.Unwrap(err)
		if err == nil {
			return NewInternalError(oldErr)
		}
	}
}

func convert(err pg.Error) Error {
	var result Error
	message := err.Field(pgMessageField)

	if strings.Contains(message, pgDuplicateErr) {
		result = NewConflictError(err)
	} else {
		result = NewInternalError(err)
	}

	return result.WithParams(err.Field(pgCodeField), err.Field(pgStatusField)).WithMessage(message)
}
