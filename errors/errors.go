package errors

type ErrorType string

const (
	Internal   = "Internal error"
	NotFound   = "Entity not found"
	Conflict   = "Entity already exists"
	BadRequest = "Found too many entities"
)

type Error interface {
	Error() string
	TypeOf(typ ErrorType) bool
	WithParams(code, status string) Error
	WithMessage(msg string) Error
	WithTag(tag *Tag) Error
	HasTag(tag *Tag) bool
}

type dbError struct {
	typ     ErrorType
	code    string
	status  string
	message string
	tags    []*Tag
	err     error
}

func (e dbError) Error() string {
	return e.message
}

func (e dbError) TypeOf(typ ErrorType) bool {
	return e.typ == typ
}

func (e *dbError) WithParams(code, status string) Error {
	e.code = code
	e.status = status
	return e
}

func (e *dbError) WithMessage(msg string) Error {
	e.message = msg
	return e
}

func (e *dbError) WithTag(tag *Tag) Error {
	if e.tags == nil {
		e.tags = make([]*Tag, 0, 1)
	}

	e.tags = append(e.tags, tag)
	return e
}

func (e *dbError) HasTag(tag *Tag) bool {
	for i := 0; i < len(e.tags); i++ {
		if e.tags[i] == tag {
			return true
		}
	}
	return false
}

func NewInternalError(err error) Error {
	return &dbError{typ: Internal, err: err}
}

func NewNotFoundError(err error) Error {
	return &dbError{typ: NotFound, err: err}
}

func NewBadRequestError(err error) Error {
	return &dbError{typ: BadRequest, err: err}
}

func NewConflictError(err error) Error {
	return &dbError{typ: Conflict, err: err}
}

func IsInternal(err error) bool {
	v, ok := err.(Error)
	if !ok {
		return true
	}
	return v.TypeOf(Internal)
}

func IsNotFound(err error) bool {
	v, ok := err.(Error)
	if !ok {
		return false
	}
	return v.TypeOf(NotFound)
}

func IsBadRequest(err error) bool {
	v, ok := err.(Error)
	if !ok {
		return false
	}
	return v.TypeOf(BadRequest)
}

func IsConflict(err error) bool {
	v, ok := err.(Error)
	if !ok {
		return false
	}
	return v.TypeOf(Conflict)
}
