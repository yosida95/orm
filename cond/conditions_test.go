package cond

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type ConditionTestCase struct {
	cond     Condition
	expected *WHERE
}

var (
	TestCases = []ConditionTestCase{
		ConditionTestCase{
			AND(Equal("Name", "Akari"), LessThanOrEqual("Age", 15)),
			&WHERE{
				"((Name = ?) AND (Age <= ?))",
				[]interface{}{
					"Akari",
					15,
				},
			},
		},
		ConditionTestCase{
			AND(OR(Equal("Name", "Akari"), Equal("Name", "Yui")), LessThanOrEqual("Age", 15)),
			&WHERE{
				"(((Name = ?) OR (Name = ?)) AND (Age <= ?))",
				[]interface{}{
					"Akari",
					"Yui",
					15,
				},
			},
		},
		ConditionTestCase{
			AND(IN("Name", "Akari", "Yui"), LessThanOrEqual("Age", 15)),
			&WHERE{
				"((Name IN (?, ?)) AND (Age <= ?))",
				[]interface{}{
					"Akari",
					"Yui",
					15,
				},
			},
		},
		ConditionTestCase{
			AND(NotEqual("Club", "Gorakubu"), LessThanOrEqual("Age", 15)),
			&WHERE{
				"((Club <> ?) AND (Age <= ?))",
				[]interface{}{
					"Gorakubu",
					15,
				},
			},
		},
	}
)

func TestCondition(t *testing.T) {
	for _, testCase := range TestCases {
		assert.Equal(t, testCase.cond.WHERE(), testCase.expected)
	}
}
