package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
