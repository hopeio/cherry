package gorm

import (
	dbi "github.com/hopeio/cherry/utils/dao/database"
	"gorm.io/gorm"
	"strings"
)

type FilterExprs dbi.FilterExprs

func (f FilterExprs) Build(odb *gorm.DB) *gorm.DB {
	var scopes []func(db *gorm.DB) *gorm.DB
	for _, filter := range f {
		filter.Field = strings.TrimSpace(filter.Field)

		if filter.Field == "" || filter.Operation == 0 || len(filter.Value) == 0 {
			continue
		}

		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where(filter.Field+" "+filter.Operation.SQL(), filter.Value...)
		})
	}
	return odb.Scopes(scopes...)
}
