package errors

type Tag struct{}

func NewTag() *Tag {
	return &Tag{}
}

func (t Tag) IsTagged(err error) bool {
	v, ok := err.(Error)
	if !ok {
		return false
	}

	return v.HasTag(&t)
}
