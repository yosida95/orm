package orm

import (
	"errors"
)

var (
	ErrUnsuitableModel        = errors.New("unsuitable model")
	ErrDestIsNotSlice         = errors.New("destination is not a slice")
	ErrTransactionHasClosed   = errors.New("transaction has closed")
	ErrInvalidQuery           = errors.New("invalid query")
)
