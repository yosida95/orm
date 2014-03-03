package orm

import (
	"database/sql"
	// "github.com/yosida95/orm/cond"
	"reflect"
)

type Client struct {
	conn   *sql.DB
	tables map[reflect.Type]*table
}

func NewClient(conn *sql.DB) *Client {
	return &Client{
		conn:   conn,
		tables: make(map[reflect.Type]*table),
	}
}

func (c *Client) initializeRecord(modelV reflect.Value) (err error) {
	if modelV.Kind() != reflect.Struct {
		err = ErrUnsuitableModel
		return
	}
	modelT := modelV.Type()

	found := false
	for i := 0; i < modelT.NumField(); i++ {
		fF := modelT.Field(i)
		if fF.Type == reflect.TypeOf(Record{}) {
			found = true
			fV := modelV.Field(i)

			rec := &Record{}
			recV := reflect.ValueOf(rec).Elem()
			recV.Set(fV)
			if !rec.initialized() {
				tbl, ok := c.tables[modelT]
				if !ok {
					tbl, err = analyzeTable(modelT, fF)
					if err != nil {
						return
					}
					c.tables[modelT] = tbl
				}
				rec.record = &record{
					client: c,
					table:  tbl,
				}
			}

			rec.parentV = modelV
			fV.Set(recV)
		}
	}

	if !found {
		err = ErrUnsuitableModel
	}
	return
}

func (c *Client) Save(r Model) (err error) {
	rV := reflect.Indirect(reflect.ValueOf(r))
	if rV.Kind() != reflect.Struct {
		err = ErrUnsuitableRecord
		return
	}

	err = c.initializeRecord(rV)
	if err != nil {
		return
	}
	return r.Save()
}

func (c *Client) QueryOne(dest Model) *Query {
	rPtrV := reflect.ValueOf(dest)
	if rPtrV.Kind() != reflect.Ptr {
		panic(ErrUnsuitableModel)
	}

	rV := rPtrV.Elem()
	err := c.initializeRecord(rV)

	rT := rV.Type()
	table := c.tables[rT]

	return &Query{
		err:     err,
		client:  c,
		dest:    rPtrV,
		isSlice: false,
		table:   table,
		limit:   1,
		offset:  0,
	}
}

func (c *Client) Query(dest []Model) *Query {
	rSliceV := reflect.ValueOf(dest)
	if rSliceV.Kind() != reflect.Slice || rSliceV.Kind() != reflect.Array {
		panic(ErrUnsuitableModel)
	}
	rSliceT := rSliceV.Type()

	var err error
	rT := rSliceT.Elem()
	table, ok := c.tables[rT]
	if !ok {
		err = c.initializeRecord(reflect.New(rT).Elem())
		if err == nil {
			table = c.tables[rT]
		}
	}

	return &Query{
		err:     err,
		client:  c,
		dest:    rSliceV,
		isSlice: true,
		table:   table,
		limit:   0,
		offset:  0,
	}
}
