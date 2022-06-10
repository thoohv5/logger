package log

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/thoohv5/logger"
	"github.com/thoohv5/logger/impl/base"
)

type entity struct {
	*base.Base
	data   logger.IData
	config *logger.Config

	logger *log.Logger

	curLevel logger.Level
}

func New(config *logger.Config) logger.ILogger {
	e := &entity{
		Base:   base.New(),
		config: config,
		data:   logger.InitData(),
	}

	clog := log.New(e.getWriter(), "", log.Lmsgprefix)

	e.logger = clog

	return e
}

func (e *entity) base(level logger.Level) *entity {
	e.curLevel = level
	return e.
		append(logger.LevelTag, base.LevelInfo(level)).
		append(logger.CallerTag, base.GetCallerInfo(3)).
		append(logger.TimeTag, time.Now().Format("2006-01-02T15:04:05.000Z0700"))
}

func (e *entity) Debug(msg string, values ...interface{}) {
	e.base(logger.Debug).append(logger.MsgTag, msg).append(logger.ParamTag, values...).output()
}

func (e *entity) Info(msg string, values ...interface{}) {
	e.base(logger.Info).append(logger.MsgTag, msg).append(logger.ParamTag, values...).output()
}

func (e *entity) Warn(msg string, values ...interface{}) {
	e.base(logger.Warn).append(logger.MsgTag, msg).append(logger.ParamTag, values...).output()
}

func (e *entity) Error(msg string, values ...interface{}) {
	e.base(logger.Error).append(logger.MsgTag, msg).append(logger.ParamTag, values...).output()
}

func (e *entity) display() string {
	return func() string {
		s := e.data.Marshal()
		e.data = logger.InitData()
		return s
	}()
}

func (e *entity) output() {
	if e.curLevel > base.Level(e.config.GetConfig().Level) {
		_ = e.logger.Output(2, e.display())
	}
}

func (e *entity) append(key string, values ...interface{}) *entity {
	e.data = e.data.Append(key, values...)
	return e
}

func (e *entity) getWriter() io.Writer {
	output := strings.Split(e.config.GetConfig().Out, ",")
	if len(output) == 0 {
		output = append(output, "std")
	}

	otw := make([]io.Writer, 0)
	for _, ot := range output {

		fc := e.config.GetFileConfig()

		var w io.Writer
		switch ot {
		case "file":
			// w, _ = os.OpenFile(fmt.Sprintf("%s/%s.log", e.path, e.currentDate()), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
			w = e.GetFileWriter(fc)
		case "std":
			fallthrough
		default:
			w = os.Stdout
		}
		otw = append(otw, w)
	}
	if len(otw) > 1 {
		return io.MultiWriter(otw...)
	}
	return otw[0]
}
