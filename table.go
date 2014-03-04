package orm

import (
	"fmt"
	"reflect"
	"strings"
)

type operation int

const (
	op_select operation = iota
	op_update operation = iota
	op_delete operation = iota
)

type column struct {
	name    string
	columnT reflect.Type
}

type table struct {
	name    string
	columns []column
}

func (t *table) query(op operation) string {
	switch op {
	case op_select:
		columns := make([]string, 0, len(t.columns))
		for _, clm := range t.columns {
			columns = append(columns, clm.name)
		}

		return fmt.Sprintf(
			strings.Join([]string{
				"SELECT",
				"    %s",
				"FROM",
				"    %s",
			}, "\n"),
			strings.Join(columns, ", "), t.name,
		)
	case op_update:
		return fmt.Sprintf("UPDATE %s", t.name)
	case op_delete:
		return fmt.Sprintf("DELETE FROM %s", t.name)
	}

	panic("unknown operation")
}

func analyzeTable(modelT reflect.Type) (tbl *table, err error) {
	tbl = &table{
		name:    "",
		columns: make([]column, 0, modelT.NumField()),
	}

	recordT := reflect.TypeOf(Record{})
	for i := 0; i < modelT.NumField(); i++ {
		fieldF := modelT.Field(i)
		fieldT := fieldF.Type

		if fieldT == recordT {
			name := fieldF.Tag.Get("orm")
			if name == "" {
				name = modelT.Name()
			}
			tbl.name = name
		} else {
			clm, ok := analyzeColumn(fieldF)
			if !ok {
				continue
			}
			tbl.columns = append(tbl.columns, clm)
		}
	}

	if tbl.name == "" {
		err = ErrUnsuitableModel
	}
	return
}

func analyzeColumn(fieldF reflect.StructField) (clm column, ok bool) {
	fieldT := fieldF.Type
	if fieldT.Kind() == reflect.Ptr {
		fieldT = fieldT.Elem()
	}

	switch fieldT.Kind() {
	case reflect.Bool:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.String:
	default:
		return
	}

	name := fieldF.Tag.Get("orm")
	if name == "-" {
		return
	} else if name == "" {
		name = fieldF.Name
	}

	clm = column{
		name:    name,
		columnT: fieldT,
	}
	ok = true
	return
}
