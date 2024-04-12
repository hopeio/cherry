package log

import (
	"fmt"
	"github.com/hopeio/cherry/utils/log/output"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	output.RegisterSink()
	SetDefaultLogger(&Config{Development: true, Level: zapcore.DebugLevel})
}

var (
	Default     *Logger
	skipLoggers = make([]*Logger, 10)
)

func SetDefaultLogger(lf *Config) {
	Default = lf.NewLogger()
}

func GetSkipLogger(skip int) *Logger {
	if skip > 10 {
		panic("skip最大不超过10")
	}
	if skipLoggers[skip] == nil {
		skipLoggers[skip] = Default.AddSkip(skip)
	}
	return skipLoggers[skip]
}

func Sync() error {
	return Default.Sync()
}

func Print(args ...interface{}) {
	if ce := Default.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Debug(args ...interface{}) {
	if ce := Default.Check(zap.DebugLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Info(args ...interface{}) {
	if ce := Default.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Warn(args ...interface{}) {
	if ce := Default.Check(zap.WarnLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Error(args ...interface{}) {
	if ce := Default.Check(zap.ErrorLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Panic(args ...interface{}) {
	if ce := Default.Check(zap.PanicLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Fatal(args ...interface{}) {
	if ce := Default.Check(zap.FatalLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Printf(template string, args ...interface{}) {
	if ce := Default.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Debugf(template string, args ...interface{}) {
	if ce := Default.Check(zap.DebugLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Infof(template string, args ...interface{}) {
	if ce := Default.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Warnf(template string, args ...interface{}) {
	if ce := Default.Check(zap.WarnLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Errorf(template string, args ...interface{}) {
	if ce := Default.Check(zap.ErrorLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Panicf(template string, args ...interface{}) {
	if ce := Default.Check(zap.PanicLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Fatalf(template string, args ...interface{}) {
	if ce := Default.Check(zap.FatalLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Debugw(msg string, fields ...zap.Field) {
	if ce := Default.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func Infow(msg string, fields ...zap.Field) {
	if ce := Default.Check(zap.InfoLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func Warnw(msg string, fields ...zap.Field) {
	if ce := Default.Check(zap.WarnLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func Errorw(msg string, fields ...zap.Field) {
	if ce := Default.Check(zap.ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}
