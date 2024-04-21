package datatypes

import (
	dbi "github.com/hopeio/cherry/utils/dao/db"
	"github.com/hopeio/cherry/utils/dao/db/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JsonT[T any] datatypes.JsonT[T]

func (*JsonT[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case dbi.Sqlite:
		return "json"
	case dbi.Mysql:
		return "json"
	case dbi.Postgres:
		return "jsonb"
	}
	return ""
}
