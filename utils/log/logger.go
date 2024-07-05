package log

import (
	"fmt"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

// Named adds a sub-scope to the logger's name. See Logger.Named for details.
func (l *Logger) Named(name string) *Logger {
	return &Logger{l.Logger.Named(name)}
}

// WithOptions warp the zap WithOptions, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	return &Logger{l.Logger.WithOptions(opts...)}
}

// With warp the zap With. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}

// Sugar warp the zap Sugar.
func (l *Logger) ZapLogger() *zap.Logger {
	return l.Logger
}

// Sugar warp the zap Sugar.
func (l *Logger) Sugar() *zap.SugaredLogger {
	l.WithOptions(zap.AddCallerSkip(-1))
	return l.Logger.Sugar()
}

// AddCore warp the zap AddCore.
func (l *Logger) AddCore(newCore zapcore.Core) *Logger {
	return l.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, newCore)
	}))
}

// AddSkip warp the zap AddCallerSkip.
func (l *Logger) AddSkip(skip int) *Logger {
	return &Logger{l.Logger.WithOptions(zap.AddCallerSkip(skip))}
}

func (l *Logger) Printf(template string, args ...any) {
	if ce := l.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// 兼容gormv1
func (l *Logger) Print(args ...any) {
	if ce := l.Check(zap.InfoLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *Logger) Debug(args ...any) {
	if ce := l.Check(zap.DebugLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

// Info uses fmt.Sprint to construct and log a message.
func (l *Logger) Info(args ...any) {
	if ce := l.Check(zap.InfoLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *Logger) Warn(args ...any) {
	if ce := l.Check(zap.WarnLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

// Error uses fmt.Sprint to construct and log a message.
func (l *Logger) Error(args ...any) {
	if ce := l.Check(zap.ErrorLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *Logger) DPanic(args ...any) {
	if ce := l.Check(zap.DPanicLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *Logger) Panic(args ...any) {
	if ce := l.Check(zap.PanicLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *Logger) Fatal(args ...any) {
	if ce := l.Check(zap.FatalLevel, trimLineBreak(fmt.Sprintln(args...))); ce != nil {
		ce.Write()
	}
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debugw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Infow(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.InfoLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Warnw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.WarnLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Errorw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l *Logger) DPanicw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.DPanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) Panicw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.PanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l *Logger) Fatalw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.FatalLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *Logger) Debugf(template string, args ...any) {
	if ce := l.Check(zap.DebugLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *Logger) Infof(template string, args ...any) {
	if ce := l.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Logger) Warnf(template string, args ...any) {
	if ce := l.Check(zap.WarnLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *Logger) Errorf(template string, args ...any) {
	if ce := l.Check(zap.ErrorLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *Logger) DPanicf(template string, args ...any) {
	if ce := l.Check(zap.DPanicLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Logger) Panicf(template string, args ...any) {
	if ce := l.Check(zap.PanicLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Logger) Fatalf(template string, args ...any) {
	if ce := l.Check(zap.FatalLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Println(args ...any) {
	if ce := l.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Debugfw(template string, args ...any) func(...zapcore.Field) {
	return func(fields ...zapcore.Field) {
		if ce := l.Check(zap.DebugLevel, fmt.Sprintf(template, args...)); ce != nil {
			ce.Write(fields...)
		}
	}
}

func (l *Logger) Infofw(template string, args ...any) func(...zapcore.Field) {
	return func(fields ...zapcore.Field) {
		if ce := l.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
			ce.Write(fields...)
		}
	}
}

func (l *Logger) Debugsw(msg string, args ...any) {
	if ce := l.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(l.sweetenFields(args)...)
	}
}

// sugar
const (
	_oddNumberErrMsg    = "Ignored key without a value."
	_nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
)

func (l *Logger) sweetenFields(args []any) []zap.Field {
	if len(args) == 0 {
		return nil
	}

	// Allocate enough space for the worst case; if users pass only structured
	// fields, we shouldn't penalize them with extra allocations.
	fields := make([]zap.Field, 0, len(args))
	var invalid invalidPairs

	for i := 0; i < len(args); {
		// This is a strongly-typed field. Consume it and move on.
		if f, ok := args[i].(zap.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		// Make sure this element isn't a dangling key.
		if i == len(args)-1 {
			l.DPanic(_oddNumberErrMsg, zap.Any(FieldIgnored, args[i]))
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string, add this pair to the slice of invalid pairs.
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			// Subsequent errors are likely, so allocate once up front.
			if cap(invalid) == 0 {
				invalid = make(invalidPairs, 0, len(args)/2)
			}
			invalid = append(invalid, invalidPair{i, key, val})
		} else {
			fields = append(fields, zap.Any(keyStr, val))
		}
		i += 2
	}

	// If we encountered any invalid key-value pairs, log an error.
	if len(invalid) > 0 {
		l.DPanic(_nonStringKeyErrMsg, zap.Array(FieldInvalid, invalid))
	}
	return fields
}

type invalidPair struct {
	position   int
	key, value any
}

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64(FieldPosition, int64(p.position))
	zap.Any(FieldKey, p.key).AddTo(enc)
	zap.Any(FieldValue, p.value).AddTo(enc)
	return nil
}

type invalidPairs []invalidPair

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	for i := range ps {
		err = multierr.Append(err, enc.AppendObject(ps[i]))
	}
	return err
}

type StdOutLevel zapcore.Level

func (l StdOutLevel) Enabled(lvl zapcore.Level) bool {
	return lvl >= zapcore.Level(l) && lvl < zapcore.ErrorLevel
}

type StdErrLevel zapcore.Level

func (l StdErrLevel) Enabled(lvl zapcore.Level) bool {
	return lvl >= zapcore.Level(l) && lvl >= zapcore.ErrorLevel
}
