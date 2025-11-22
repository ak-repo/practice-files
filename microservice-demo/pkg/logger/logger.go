package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init() error {

	if Log != nil {
		return nil
	}

	l, err := zap.NewProduction()
	if err != nil {
		return nil
	}

	Log = l
	return nil
}
