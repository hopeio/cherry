package gormi

import (
	"context"
	contexti "github.com/hopeio/cherry/utils/context"
	"gorm.io/gorm"
)

func NewTraceDB(db *gorm.DB, ctx context.Context, traceId string) *gorm.DB {
	return db.Session(&gorm.Session{Context: contexti.SetTranceId(ctx, traceId), NewDB: true})
}
