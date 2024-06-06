package gorm

import (
	"context"
	loggeri "github.com/hopeio/cherry/utils/dao/database/gorm/logger"
	"github.com/hopeio/cherry/utils/log"
	contexti "github.com/hopeio/cherry/utils/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBWithLogger(db *gorm.DB, log *log.Logger, conf *logger.Config) *gorm.DB {
	return db.Session(&gorm.Session{
		Logger: &loggeri.Logger{Logger: log.Logger,
			Config: conf,
		}})
}

func GetDBWithContext(db *gorm.DB, ctx context.Context) *gorm.DB {
	return db.Session(&gorm.Session{Context: ctx})
}

func NewTraceDB(db *gorm.DB, ctx context.Context, traceId string) *gorm.DB {
	return db.Session(&gorm.Session{Context: contexti.SetTranceId(ctx, traceId), NewDB: true})
}
