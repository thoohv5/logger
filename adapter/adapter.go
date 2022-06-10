package logger

import (
	"github.com/thoohv5/logger"
	"github.com/thoohv5/logger/impl/log"
	"github.com/thoohv5/logger/impl/zap"
)

type Type int

const (
	Log Type = iota
	Zap
)

func Adapter(t Type, config *logger.Config) logger.ILogger {
	var l logger.ILogger
	switch t {
	case Log:
		l = log.New(config)
	case Zap:
		l = zap.New(config)
	}
	return l
}
