package orm

import (
	"reflect"
)

type Query struct {
	err     error
	client  *Client
	dest    reflect.Value
	isSlice bool
	table   *table
	limit   int
	offset  int
}

func (q *Query) Limit(limit int) *Query {
	if limit > 1 && !q.isSlice {
		panic(ErrDestIsNotSlice)
	}

	q.limit = limit
	return q
}

func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}
