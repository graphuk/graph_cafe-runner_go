package data

import (
	"fmt"
	"reflect"

	"github.com/asdine/storm"
)

type UnitOfWork struct {
	DB storm.Node
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// command is a database function working with the database (read-write) and returning no results
func (t *UnitOfWork) ExecuteCommand(command func(tx storm.Node) (interface{}, error)) (interface{}, error) {

	var err error
	var res interface{}

	if fmt.Sprint(reflect.TypeOf(t.DB)) == `*storm.DB` { // if we are not inside a transaction - start transaction
		tx, err := t.DB.Begin(true)
		check(err)
		defer tx.Rollback()

		res, err = command(tx)
		check(tx.Commit())
	} else { // if we are inside a transaction - run func in context of exist transaction
		res, err = command(t.DB)
	}

	return res, err
}

// command is a database function working with the database (read only) and returning no results
func (t *UnitOfWork) ExecuteQuery(command func(tx storm.Node) (interface{}, error)) (interface{}, error) {
	var err error
	var res interface{}

	if fmt.Sprint(reflect.TypeOf(t.DB)) == `*storm.DB` { // if we are not inside a transaction - start transaction
		tx, err := t.DB.Begin(false)
		check(err)
		defer tx.Rollback()

		res, err = command(tx)
		check(tx.Commit())
	} else { // if we are inside a transaction - run func in context of exist transaction
		res, err = command(t.DB)
	}

	return res, err
}
