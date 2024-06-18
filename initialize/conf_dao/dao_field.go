package conf_dao

import (
	"io"
	"reflect"
)

var DaoFieldType = reflect.TypeOf((*DaoField)(nil)).Elem()

type DaoField interface {
	Config() any
	Init() error
	io.Closer
}
