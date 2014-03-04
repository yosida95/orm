package orm

import (
	"database/sql"
	"fmt"
	"reflect"
)

type Session struct {
	tx            *sql.Tx
	client        *Client
	inTransaction bool
}

func (s *Session) Save(model Model) error {
	return nil
}

func (s *Session) Commit() error {
	if !s.inTransaction {
		panic(ErrTransactionHasClosed)
	}

	s.inTransaction = false
	return s.tx.Commit()
}

func (s *Session) Rollback() error {
	if !s.inTransaction {
		panic(ErrTransactionHasClosed)
	}

	s.inTransaction = false
	return s.tx.Rollback()
}

/*
func (s *Session) exec(q *Query) error {
	if !s.inTransaction {
		panic(ErrTransactionHasClosed)
	}

	// tx.Exec
	return nil
}
*/

func (s *Session) query(q *Query) error {
	if !s.inTransaction {
		panic(ErrTransactionHasClosed)
	}

	args := make([]interface{}, 0)

	exp := q.table.query(op_select)
	if q.where != nil {
		where := q.where.WHERE()
		exp = fmt.Sprintf("%s\nWHERE\n    %s", exp, where.Expression())

		cargs := make([]interface{}, len(args))
		copy(cargs, args)
		wargs := where.Args()

		args = make([]interface{}, len(cargs)+len(wargs))
		copy(args, cargs)
		copy(args[len(cargs):], wargs)
	}

	if !q.isSlice {
		exp = fmt.Sprintf("%s\n%s", exp, "LIMIT 1")
	} else if q.limit > 0 {
		exp = fmt.Sprintf("%s\n%s", exp, fmt.Sprintf("LIMIT %d", q.limit))
	}

	if q.offset > 0 {
		exp = fmt.Sprintf("%s\n%s", exp, fmt.Sprintf("OFFSET %d", q.offset))
	}

	// s.tx.Query(exp, args...)
	println(exp)
	return nil
}

func (s *Session) Query(models interface{}) *Query {
	if !s.inTransaction {
		panic(ErrTransactionHasClosed)
	}

	destV := reflect.ValueOf(models)
	if destV.Kind() != reflect.Slice && destV.Kind() != reflect.Array {
		panic(ErrUnsuitableModel)
	}

	destT := destV.Type()
	modelT := destT.Elem()
	_, ok := reflect.New(modelT).Interface().(Model)
	if !ok {
		panic(ErrUnsuitableModel)
	}

	table, err := s.client.getTable(modelT)
	if err != nil {
		panic(err)
	}

	return s.buildQuery(table, destV)
}

func (s *Session) QueryOne(model Model) *Query {
	if !s.inTransaction {
		panic(ErrTransactionHasClosed)
	}

	destV := reflect.ValueOf(model)
	if destV.Kind() != reflect.Ptr && destV.Kind() != reflect.Struct {
		panic(ErrUnsuitableModel)
	}

	modelV := reflect.Indirect(destV)
	modelT := modelV.Type()
	table, err := s.client.getTable(modelT)
	if err != nil {
		panic(err)
	}

	return s.buildQuery(table, destV)
}

func (s *Session) buildQuery(table *table, destV reflect.Value) *Query {
	isSlice := false
	if destV.Kind() == reflect.Slice || destV.Kind() == reflect.Array {
		isSlice = true
	}

	return &Query{
		session: s,
		where:   nil,
		dest:    destV,
		isSlice: isSlice,
		table:   table,
		limit:   0,
		offset:  0,
	}
}
