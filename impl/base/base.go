package base

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/natefinch/lumberjack"

	"github.com/thoohv5/logger"
	"github.com/thoohv5/logger/util"
)

type Base struct {
	options logger.Options
}

func New(opts ...logger.Option) *Base {
	o := logger.Options{
		Skip: 3,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &Base{
		options: o,
	}
}

func (b *Base) GetFileWriter(fc *logger.File) io.Writer {
	if strings.HasPrefix(fc.Path, ".") {
		fc.Path = util.AbPath(fc.Path)
	}
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", fc.Path, fc.FileName),
		MaxSize:    fc.MaxAge,
		MaxBackups: fc.MaxBackups,
		MaxAge:     fc.MaxAge,
		Compress:   fc.Compress,
	}
}

func Level(levelStr string) logger.Level {
	var level logger.Level
	switch levelStr {
	case "debug":
		level = logger.Debug
	case "info":
		level = logger.Info
	case "warn":
		level = logger.Warn
	case "error":
		level = logger.Error
	}
	return level
}

func LevelInfo(level logger.Level) string {
	levelVal := "debug"
	switch level {
	case logger.Info:
		levelVal = "info"
	case logger.Warn:
		levelVal = "warn"
	case logger.Error:
		levelVal = "error"

	}
	return levelVal
}

type CallerInfo struct {
	Func string `json:"fun,omitempty"`
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

func GetCallerInfo(skip int) string {

	_, file, lineNo, ok := runtime.Caller(skip)
	// pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		return "runtime.Caller() failed"
	}
	// funcName := runtime.FuncForPC(pc).Name()
	// fileName := path.Base(file) // Base函数返回路径的最后一个元素
	return fmt.Sprintf("%s:%d", file, lineNo)
}

func (b *Base) GetSkip() int {
	return b.options.GetSkip()
}
