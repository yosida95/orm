package orm

import (
	"reflect"
)

type table struct {
	name string
}

func analyzeTable(modelT reflect.Type, fF reflect.StructField) (tbl *table, err error) {
	name := fF.Tag.Get("orm")
	if name == "" {
		name = modelT.Name()
	}

	tbl = &table{
		name: name,
	}
	return
}
