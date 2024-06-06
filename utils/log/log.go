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
	stackLogger = defaultLogger.WithOptions(zap.WithCaller(true), zap.AddStacktrace(zapcore.ErrorLevel))
}

type skipLogger struct {
	*Logger
	needUpdate bool
}

var (
	defaultLogger *Logger
	stackLogger   *Logger
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

	defaultLogger = lf.NewLogger(cores...)
	skipLoggers[0].Logger = defaultLogger.WithOptions(zap.WithCaller(false))
	for i := 1; i < len(skipLoggers); i++ {
		if skipLoggers[i].Logger != nil {
			skipLoggers[i].needUpdate = true
		}
	}
}

func GetCallerSkipLogger(skip int) *Logger {
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

func GetNoCallerLogger() *Logger {
	return skipLoggers[0].Logger
}

func Sync() error {
	return defaultLogger.Sync()
}

func Print(args ...any) {
	if ce := defaultLogger.Check(zap.InfoLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func Debug(args ...any) {
	if ce := defaultLogger.Check(zap.DebugLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func Info(args ...any) {
	if ce := defaultLogger.Check(zap.InfoLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func Warn(args ...any) {
	if ce := defaultLogger.Check(zap.WarnLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func Error(args ...any) {
	if ce := defaultLogger.Check(zap.ErrorLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func Panic(args ...any) {
	if ce := defaultLogger.Check(zap.PanicLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func Fatal(args ...any) {
	if ce := defaultLogger.Check(zap.FatalLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func Printf(template string, args ...any) {
	if ce := defaultLogger.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Debugf(template string, args ...any) {
	if ce := defaultLogger.Check(zap.DebugLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Infof(template string, args ...any) {
	if ce := defaultLogger.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Warnf(template string, args ...any) {
	if ce := defaultLogger.Check(zap.WarnLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Errorf(template string, args ...any) {
	if ce := defaultLogger.Check(zap.ErrorLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Panicf(template string, args ...any) {
	if ce := defaultLogger.Check(zap.PanicLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Fatalf(template string, args ...any) {
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

func Panicw(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.PanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func Fatalw(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.FatalLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func Println(args ...any) {
	if ce := defaultLogger.Check(zap.InfoLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func ErrorStack(args ...any) {
	if ce := stackLogger.Check(zap.ErrorLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func PanicStack(args ...any) {
	if ce := stackLogger.Check(zap.PanicLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func FatalStack(args ...any) {
	if ce := stackLogger.Check(zap.FatalLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

func ErrorStackf(template string, args ...any) {
	if ce := stackLogger.Check(zap.ErrorLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func FatalStackf(template string, args ...any) {
	if ce := stackLogger.Check(zap.FatalLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func ErrorStackw(msg string, fields ...zap.Field) {
	if ce := stackLogger.Check(zap.ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func PanicStackw(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.PanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func FatalStackw(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.FatalLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}
