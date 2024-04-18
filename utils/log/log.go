package log

import (
	"fmt"
	"github.com/hopeio/cherry/utils/log/output"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

func init() {
	output.RegisterSink()
	SetDefaultLogger(&Config{Development: true, Level: zapcore.DebugLevel})
}

type skipLogger struct {
	*Logger
	needUpdate bool
}

var (
	defaultLogger *Logger
	skipLoggers   = make([]skipLogger, 10)
	mu            sync.Mutex
)

//go:nosplit
func Default() *Logger {
	return defaultLogger
}

func SetDefaultLogger(lf *Config, cores ...zapcore.Core) {
	mu.Lock()
	defer mu.Unlock()
	for i := 1; i < len(skipLoggers); i++ {
		if skipLoggers[i].Logger != nil {
			skipLoggers[i].needUpdate = true
		}
	}
	defaultLogger = lf.NewLogger(cores...)
	skipLoggers[0].Logger = defaultLogger
}

func GetSkipLogger(skip int) *Logger {
	if skip > 10 {
		panic("skip最大不超过10")
	}
	if skipLoggers[skip].needUpdate || skipLoggers[skip].Logger == nil {
		mu.Lock()
		skipLoggers[skip].Logger = defaultLogger.AddSkip(skip)
		skipLoggers[skip].needUpdate = false
		mu.Unlock()
	}
	return skipLoggers[skip].Logger
}

func Sync() error {
	return defaultLogger.Sync()
}

func Print(args ...interface{}) {
	if ce := defaultLogger.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Println(args ...interface{}) {
	if ce := defaultLogger.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Debug(args ...interface{}) {
	if ce := defaultLogger.Check(zap.DebugLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Info(args ...interface{}) {
	if ce := defaultLogger.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Warn(args ...interface{}) {
	if ce := defaultLogger.Check(zap.WarnLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Error(args ...interface{}) {
	if ce := defaultLogger.Check(zap.ErrorLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Panic(args ...interface{}) {
	if ce := defaultLogger.Check(zap.PanicLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Fatal(args ...interface{}) {
	if ce := defaultLogger.Check(zap.FatalLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Printf(template string, args ...interface{}) {
	if ce := defaultLogger.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Debugf(template string, args ...interface{}) {
	if ce := defaultLogger.Check(zap.DebugLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Infof(template string, args ...interface{}) {
	if ce := defaultLogger.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Warnf(template string, args ...interface{}) {
	if ce := defaultLogger.Check(zap.WarnLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Errorf(template string, args ...interface{}) {
	if ce := defaultLogger.Check(zap.ErrorLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Panicf(template string, args ...interface{}) {
	if ce := defaultLogger.Check(zap.PanicLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Fatalf(template string, args ...interface{}) {
	if ce := defaultLogger.Check(zap.FatalLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Debugw(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func Infow(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.InfoLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func Warnw(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.WarnLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func Errorw(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}
