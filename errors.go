package orm

import (
	"errors"
)

var (
	ErrUnsuitableRecord = errors.New("unsuitable record")
	ErrUnsuitableModel  = errors.New("unsuitable model")
	ErrDestIsNotSlice   = errors.New("destination is not a slice")
)
