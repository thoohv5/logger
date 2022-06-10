package zap

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/thoohv5/logger"
	"github.com/thoohv5/logger/impl/base"
)

type entity struct {
	*base.Base
	zap    *zap.Logger
	config *logger.Config
}

const (
	OutSep = ","
)

func New(config *logger.Config) logger.ILogger {

	e := &entity{
		config: config,
	}
	caller := zap.AddCaller()

	// 设置初始化字段
	// filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	// log := zap.New(core, caller, zap.Development(), filed)
	z := zap.New(e.getWriter(), caller, zap.Development())
	e.zap = z
	return e
}

func (e *entity) Debug(msg string, values ...interface{}) {
	e.zap.Debug(msg, zap.Any(logger.ParamTag, values))
}

func (e *entity) Info(msg string, values ...interface{}) {
	e.zap.Info(msg, zap.Any(logger.ParamTag, values))
}

func (e *entity) Warn(msg string, values ...interface{}) {
	e.zap.Warn(msg, zap.Any(logger.ParamTag, values))
}

func (e *entity) Error(msg string, values ...interface{}) {
	e.zap.Error(msg, zap.Any(logger.ParamTag, values))
}

func (e *entity) getWriter() zapcore.Core {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        logger.TimeTag,
		LevelKey:       logger.LevelTag,
		NameKey:        "log",
		CallerKey:      logger.CallerTag,
		MessageKey:     logger.MsgTag,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	ws := make([]zapcore.WriteSyncer, 0)
	for _, out := range strings.Split(e.config.GetConfig().Out, OutSep) {
		switch out {
		case "std":
			ws = append(ws, zapcore.AddSync(os.Stdout))
		default:
			ws = append(ws, zapcore.AddSync(e.GetFileWriter(e.config.GetFileConfig())))
		}
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(ws...),
		zap.NewAtomicLevelAt(parseLevel(e.config.GetConfig().Level)),
	)

}

// 日志类别: debug, warn, info，error
func parseLevel(level string) zapcore.Level {
	zl := zapcore.DebugLevel
	switch level {
	case "debug":
		zl = zapcore.DebugLevel
	case "warn":
		zl = zapcore.WarnLevel
	case "info":
		zl = zapcore.InfoLevel
	case "error":
		zl = zapcore.ErrorLevel
	}
	return zl
}
