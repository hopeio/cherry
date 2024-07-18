package model

import "github.com/hopeio/utils/dao/database/datatypes"

type TestJson struct {
	ID        uint                      `json:"id" gorm:"primaryKey"`
	JsonArray datatypes.JsonTArray[Foo] `json:"json_array" gorm:"jsonb"`
}

type Foo struct {
	A int
	B string
}
