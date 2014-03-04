package orm

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"testing"
)

func TestRecord(t *testing.T) {
	rec := &Record{}
	if rec.initialized() {
		t.Log("record was not initialized, but initialized() returns true")
		t.Fail()
	}

	rec.record = &record{}
	if rec.initialized() {
		t.Log("record was not initialized, but initialized() returns true")
		t.Fail()
	}

	dbconn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	rec.record.client = NewClient(dbconn)

	if rec.initialized() {
		t.Log("record was not initialized, but initialized() returns true")
		t.Fail()
	}

	rec.record.table = &table{
		name: "dummy",
	}

	if !rec.initialized() {
		t.Log("record was initialized, but initialized() returns false")
		t.Fail()
	}
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

	/*
		conn := NewClient(dbconn)

		dummy := &DummyModel{
			Field1:   "foo",
			Field2:   "bar",
			Internal: "baz",
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
	*/
}
