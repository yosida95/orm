package orm

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"testing"
)

type DummyModel struct {
	Field1   string `orm:"field1"`
	Field2   string `orm:"field2"`
	Internal string `orm:"-"`
	Record   `orm:"dummy"`
}

func Test(t *testing.T) {
	dbconn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	_, err = dbconn.Exec(strings.Join([]string{
		"CREATE TABLE dummy (",
		"    field1 CHAR,",
		"    field2 CHAR",
		");",
	}, "\n"))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	conn := NewClient(dbconn)

	dummy := &DummyModel{
		Field1:   "foo",
		Field2:   "bar",
		Internal: "baz",
	}
	err = conn.Save(dummy)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	record := dummy.Record
	if record.record.client != conn {
		t.Logf("current value of conn is not expected one: %v", record.record.client)
		t.Fail()
	}

	if record.record.table.name != "dummy" {
		t.Logf("failed to get the table name")
		t.Fail()
	}
}
