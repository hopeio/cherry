package gormi

import (
	"context"
	"github.com/hopeio/cherry/utils/scheduler/crawler"
	"gorm.io/gorm"
)

type Cursor struct {
	Type   string `json:"type" gorm:"primaryKey"`
	Cursor string
	Prev   string
	Next   string
}

func GetCursor(db *gorm.DB, mod string) (*Cursor, error) {
	var cursor Cursor
	err := db.Where(`type = ?`, mod).First(&cursor).Error
	if err != nil {
		return nil, err
	}
	return &cursor, nil
}

func EndCallback(db *gorm.DB, mod string) {
	db.Exec(`UPDATE cursor SET prev = next WHERE typpe = '` + mod + `'`)
	db.Exec(`UPDATE cursor SET cursor = '' WHERE typpe = '` + mod + `'`)
}

func EndCallbackReq(db *gorm.DB, mod string) *crawler.Request {
	return &crawler.Request{
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			EndCallback(db, mod)
			return nil, nil
		},
	}
}
