package context

import (
	"github.com/hopeio/cherry/protobuf/errcode"
	"github.com/hopeio/cherry/utils/log"
	"go.uber.org/zap"
)

func (c *Context) Error(args ...interface{}) {
	args = append(args, log.FieldTraceId, c.TraceID)
	log.Error(args...)
}

func (c *Context) HandleError(err error) {
	if err != nil {
		log.Errorw(err.Error(), zap.String(log.FieldTraceId, c.TraceID))
	}
}

func (c *Context) ErrorLog(err, originErr error, funcName string) error {
	// caller 用原始logger skip刚好
	log.GetCallerSkipLogger(1).Errorw(originErr.Error(), zap.String(log.FieldTraceId, c.TraceID), zap.Int(log.FieldType, errcode.Code(err)), zap.String(log.FieldPosition, funcName))
	return err
}
