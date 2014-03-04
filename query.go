package orm

import (
	"github.com/yosida95/orm/cond"
	"reflect"
)

type Query struct {
	session *Session
	where   cond.Condition
	dest    reflect.Value
	isSlice bool
	table   *table
	limit   int
	offset  int
}

func (q *Query) Where(where cond.Condition) *Query {
	if q.where == nil {
		q.where = where
	} else {
		q.where = cond.AND(q.where, where)
	}

	return q
}

func (q *Query) Limit(limit int) *Query {
	if limit > 1 && !q.isSlice {
		panic(ErrDestIsNotSlice)
	}

	q.limit = limit
	return q
}

func (q *Query) Offset(offset int) *Query {
	if offset < 0 {
		panic(ErrInvalidQuery)
	}

	q.offset = offset
	return q
}

func (q *Query) Do() error {
	return q.session.query(q)
}
