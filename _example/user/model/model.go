package model

import "github.com/hopeio/cherry/utils/dao/db/datatypes"

type TestJson struct {
	ID        uint                      `json:"id" gorm:"primaryKey"`
	JsonArray datatypes.JsonArrayT[Foo] `json:"json_array" gorm:"jsonb"`
}

type Foo struct {
	A int
	B string
}