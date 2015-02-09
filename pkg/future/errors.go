package future

import (
	"errors"
)

var (
	DiscardedBeforeDoneError = errors.New("future discarded before done")
	NilFutureError           = errors.New("nil future")
	nilTypeError             = newTypeError("expected non-nil object, but got <nil> invalid kind")
)

type TypeError string

func (t *TypeError) Error() string {
	if t == nil {
		return ""
	}
	return string(*t)
}

func newTypeError(s string) *TypeError {
	msg := TypeError(s)
	return &msg
}
