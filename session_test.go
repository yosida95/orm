package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/yosida95/orm/cond"
	"reflect"
	"testing"
)

type DummyModel struct {
	Field1   string `orm:"field1"`
	Field2   string `orm:"field2"`
	Internal string `orm:"-"`
	Record   `orm:"dummy"`
}

func TestSession(t *testing.T) {
	dbconn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	client := NewClient(dbconn)

	session, err := client.Session()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if session.client != client {
		t.FailNow()
	}
	if !session.inTransaction {
		t.FailNow()
	}

	models := make([]DummyModel, 0)
	query := session.Query(models)
	assert.Equal(t, query.session, session)
	assert.Equal(t, query.where, nil)
	assert.Equal(t, query.dest.Type(), reflect.TypeOf(models))
	assert.Equal(t, query.isSlice, true)
	assert.Equal(t, query.limit, 0)
	assert.Equal(t, query.offset, 0)
	query.Do()

	query.Where(cond.AND(
		cond.Equal("field1", "value"),
		cond.Equal("field2", "value"),
	))
	query.Do()
	table1 := query.table

	model := new(DummyModel)
	query = session.QueryOne(model)
	assert.Equal(t, query.session, session)
	assert.Equal(t, query.where, nil)
	assert.Equal(t, query.dest.Type(), reflect.TypeOf(model))
	assert.Equal(t, query.isSlice, false)
	assert.Equal(t, query.limit, 0)
	assert.Equal(t, query.offset, 0)
	query.Do()
	table2 := query.table

	assert.Equal(t, table1, table2)
	assert.Equal(t, table1.name, "dummy")

	modelT := reflect.TypeOf(model).Elem()
	field1F, _ := modelT.FieldByName("Field1")
	field2F, _ := modelT.FieldByName("Field2")
	columns := []column{
		column{
			name:    "field1",
			columnT: field1F.Type,
		},
		column{
			name:    "field2",
			columnT: field2F.Type,
		},
	}
	assert.Equal(t, table1.columns, columns)
}
