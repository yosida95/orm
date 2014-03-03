package orm

import (
	"reflect"
)

type record struct {
	client *Client
	table  *table
}

func (r *record) initialized() bool {
	if r.client == nil {
		return false
	} else if r.table == nil {
		return false
	}

	return true
}

type Model interface {
	Save() error
}

type Record struct {
	record  *record
	parentV reflect.Value
}

func (r *Record) initialized() bool {
	if r.record == nil || !r.record.initialized() {
		return false
	}

	return true
}

func (r *Record) Save() (err error) {
	if r.record == nil {
		/*
			r.record = &record{
				client: client,
			}
		*/
	} else {
		if r.record.client == nil {
			err = ErrUnsuitableRecord
			return
		}
	}

	return
}
