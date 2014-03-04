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

func (c *Client) Session() (s *Session, err error) {
	tx, err := c.conn.Begin()
	if err != nil {
		return
	}

	s = &Session{
		tx:            tx,
		client:        c,
		inTransaction: true,
	}
	return
}

func (c *Client) getTable(modelT reflect.Type) (tbl *table, err error) {
	tbl, ok := c.tables[modelT]
	if ok {
		return
	}

	tbl, err = analyzeTable(modelT)
	if err != nil {
		return
	}
	c.tables[modelT] = tbl

	return
}
