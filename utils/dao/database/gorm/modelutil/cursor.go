package modelutil

import (
	"github.com/hopeio/cherry/utils/dao/database/modelutil"
	"github.com/hopeio/cherry/utils/types/model"
	"gorm.io/gorm"
)

func GetCursor(db *gorm.DB, typ string) (*model.Cursor, error) {
	var cursor model.Cursor
	err := db.Where(`type = ?`, typ).First(&cursor).Error
	if err != nil {
		return nil, err
	}
	return &cursor, nil
}

func EndCallback(db *gorm.DB, typ string) {
	db.Exec(modelutil.EndCallbackSQL(typ))
}
