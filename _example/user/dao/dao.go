package dao

import (
	"github.com/hopeio/cherry/context/http_context"
	"github.com/hopeio/cherry/utils/log"
)

type userDao struct {
	*http_context.Context
}

func GetDao(ctx *http_context.Context) *userDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &userDao{ctx}
}
