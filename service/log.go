package service

import (
	"bytes"
	"fmt"
	"runtime"

	logger "github.com/GXK666/eosTransfer/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

var log *zap.SugaredLogger

func setupLog() {
	sugar := logger.Plain.WithOptions(zap.AddCallerSkip(2)).Sugar()
	log = sugar.Named("rpc")

	grpclog.SetLoggerV2(&grpclogger{
		level: viper.GetInt("rpc.logLevel"),
	})
}

func LogPanicHandler(r interface{}) {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("Recovered from panic: %#v (%v)", r, r))
	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		buffer.WriteString(fmt.Sprintf("    %d: %v:%v\n", i, file, line))
	}
	log.Error(buffer.String())
}

type grpclogger struct {
	level int
}

func (g *grpclogger) Info(args ...interface{}) {
	log.Info(args...)
}

func (g *grpclogger) Infoln(args ...interface{}) {
	log.Info(args...)
}

func (g *grpclogger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func (g *grpclogger) Warning(args ...interface{}) {
	log.Warn(args...)
}

func (g *grpclogger) Warningln(args ...interface{}) {
	log.Warn(args...)
}

func (g *grpclogger) Warningf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func (g *grpclogger) Error(args ...interface{}) {
	log.Error(args...)
}

func (g *grpclogger) Errorln(args ...interface{}) {
	log.Error(args)
}

func (g *grpclogger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func (g *grpclogger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func (g *grpclogger) Fatalln(args ...interface{}) {
	log.Fatal(args...)
}

func (g *grpclogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func (g *grpclogger) V(l int) bool {
	return l >= g.level
}
