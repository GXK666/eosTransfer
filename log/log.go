package log

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Plain *zap.Logger
var Sugar *zap.SugaredLogger

func Init() {
	cfg := viper.Sub("log")
	var level zapcore.Level
	err := level.Set(cfg.GetString("level"))
	if err != nil {
		panic(err.Error())
	}

	var conf zap.Config
	if level == zapcore.DebugLevel {
		conf = zap.NewDevelopmentConfig()
	} else {
		conf = zap.NewProductionConfig()
	}
	conf.Level.SetLevel(level)
	conf.DisableStacktrace = cfg.GetBool("disableStacktrace")

	Plain, err = conf.Build()
	if err != nil {
		panic(err.Error())
	}
	Sugar = Plain.Sugar()
}

type SugaredLogger struct {
	*zap.SugaredLogger
}

func Named(name string) *SugaredLogger {
	return &SugaredLogger{Sugar.Named(name)}
}

func (l *SugaredLogger) Event(event string) *zap.SugaredLogger {
	return l.With("event", event)
}

func (l *SugaredLogger) Context(context string) *zap.SugaredLogger {
	return l.With("context", context)
}

func Debug(args ...interface{}) {
	Sugar.Debug(args...)
}

func Info(args ...interface{}) {
	Sugar.Info(args...)
}

func Warn(args ...interface{}) {
	Sugar.Warn(args...)
}

func Error(args ...interface{}) {
	Sugar.Error(args...)
}

func DPanic(args ...interface{}) {
	Sugar.DPanic(args...)
}

func Panic(args ...interface{}) {
	Sugar.Panic(args...)
}

func Fatal(args ...interface{}) {
	Sugar.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	Sugar.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	Sugar.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Sugar.Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	Sugar.DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	Sugar.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Sugar.Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	Sugar.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	Sugar.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	Sugar.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	Sugar.Errorw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	Sugar.DPanicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	Sugar.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	Sugar.Fatalw(msg, keysAndValues...)
}

func Sync() error {
	return Sugar.Sync()
}
