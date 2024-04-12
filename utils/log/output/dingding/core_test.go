package dingding

import (
	"github.com/hopeio/cherry/utils/log"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestDingDing(t *testing.T) {

	log.Default = (&log.Config{
		Development: false,
		Caller:      true,
		Level:       zapcore.DebugLevel,
		OutputPaths: log.OutPutPaths{},
		ModuleName:  "",
	}).NewLogger(NewCore("", "", zapcore.DebugLevel))
	log.Info("测试")
}
