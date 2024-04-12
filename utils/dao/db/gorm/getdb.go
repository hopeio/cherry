package gormi

import (
	loggeri "github.com/hopeio/cherry/utils/dao/db/gorm/logger"
	"github.com/hopeio/cherry/utils/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDB(db *gorm.DB, log *log.Logger, conf *logger.Config) *gorm.DB {
	return db.Session(&gorm.Session{
		Logger: &loggeri.Logger{Logger: log.Logger,
			Config: conf,
		}})
}
