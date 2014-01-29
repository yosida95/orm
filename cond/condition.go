package cond

import (
	"fmt"
	"strings"
)

type WHERE struct {
	exp  string
	args []interface{}
}

func (sql *WHERE) Expression() string {
	return sql.exp
}

func (sql *WHERE) Args() []interface{} {
	return sql.args
}

type Condition interface {
	WHERE() *WHERE
}

func AND(left, right Condition) Condition {
	return &and{
		Left:  left,
		Right: right,
	}
}

type and struct {
	Left  Condition
	Right Condition
}

func (cond *and) WHERE() *WHERE {
	left := cond.Left.WHERE()
	right := cond.Right.WHERE()

	largs := left.Args()
	rargs := right.Args()

	args := make([]interface{}, len(largs)+len(rargs))
	copy(args, largs)
	copy(args[len(largs):], rargs)

	return &WHERE{
		fmt.Sprintf("(%s AND %s)", left.Expression(), right.Expression()),
		args,
	}
}

func OR(left, right Condition) Condition {
	return &or{
		Left:  left,
		Right: right,
	}
}

type or struct {
	Left  Condition
	Right Condition
}

func (cond *or) WHERE() *WHERE {
	left := cond.Left.WHERE()
	right := cond.Right.WHERE()

	largs := left.Args()
	rargs := right.Args()

	args := make([]interface{}, len(largs)+len(rargs))
	copy(args, largs)
	copy(args[len(largs):], rargs)

	return &WHERE{
		fmt.Sprintf("(%s OR %s)", left.Expression(), right.Expression()),
		args,
	}
}

func Equal(column string, value interface{}) Condition {
	return &equal{
		column,
		value,
	}
}

type equal struct {
	Column string
	Value  interface{}
}

func (cond *equal) WHERE() *WHERE {
	return &WHERE{
		fmt.Sprintf("(%s = ?)", cond.Column),
		[]interface{}{cond.Value},
	}
}

func NotEqual(column string, value interface{}) Condition {
	return &notEqual{
		column,
		value,
	}
}

type notEqual struct {
	Column string
	Value  interface{}
}

func (cond *notEqual) WHERE() *WHERE {
	return &WHERE{
		fmt.Sprintf("(%s <> ?)", cond.Column),
		[]interface{}{cond.Value},
	}
}

func LessThan(column string, value interface{}) Condition {
	return &lessThan{
		column,
		value,
	}
}

type lessThan struct {
	Column string
	Value  interface{}
}

func (cond *lessThan) WHERE() *WHERE {
	return &WHERE{
		fmt.Sprintf("(%s < ?)", cond.Column),
		[]interface{}{cond.Value},
	}
}

func GreaterThan(column string, value interface{}) Condition {
	return &greaterThan{
		column,
		value,
	}
}

type greaterThan struct {
	Column string
	Value  interface{}
}

func (cond *greaterThan) WHERE() *WHERE {
	return &WHERE{
		fmt.Sprintf("(%s > ?)", cond.Column),
		[]interface{}{cond.Value},
	}
}

func LessThanOrEqual(column string, value interface{}) Condition {
	return &lessThanOrEqual{
		column,
		value,
	}
}

type lessThanOrEqual struct {
	Column string
	Value  interface{}
}

func (cond *lessThanOrEqual) WHERE() *WHERE {
	return &WHERE{
		fmt.Sprintf("(%s <= ?)", cond.Column),
		[]interface{}{cond.Value},
	}
}

func GreaterThanOrEqual(column string, value interface{}) Condition {
	return &greaterThanOrEqual{
		column,
		value,
	}
}

type greaterThanOrEqual struct {
	Column string
	Value  interface{}
}

func (cond *greaterThanOrEqual) WHERE() *WHERE {
	return &WHERE{
		fmt.Sprintf("(%s >= ?)", cond.Column),
		[]interface{}{cond.Value},
	}
}

func IN(column string, values ...interface{}) Condition {
	return &in{
		column,
		values,
	}
}

type in struct {
	column string
	values []interface{}
}

func (cond *in) WHERE() *WHERE {
	values := make([]string, 0, len(cond.values))
	for i := 0; i < len(cond.values); i++ {
		values = append(values, "?")
	}

	return &WHERE{
		fmt.Sprintf("(%s IN (%s))", cond.column, strings.Join(values, ", ")),
		cond.values,
	}
}
