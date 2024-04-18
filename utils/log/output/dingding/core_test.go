package dingding

import (
	"github.com/hopeio/cherry/utils/log"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestDingDing(t *testing.T) {

	log.SetDefaultLogger(&log.Config{
		Development: false,
		Caller:      true,
		Level:       zapcore.DebugLevel,
		OutputPaths: log.OutPutPaths{},
		AppName:     "",
	}, NewCore("", "", zapcore.DebugLevel))
	log.Info("测试")
}
