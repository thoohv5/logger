package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/thoohv5/logger"
	"github.com/thoohv5/logger/impl/base"
)

type entity struct {
	*base.Base

	w io.Writer

	ch     chan logger.IData
	msgBuf []logger.IData
	cf     *logger.Config
	buf    *bytes.Buffer
}

const (
	MsgChanLen        = 10000
	MaxMsgBufLen      = 100
	FlushTimeInterval = time.Second
)

func New(config *logger.Config, opts ...logger.Option) logger.ILogger {

	e := &entity{
		Base:   base.New(opts...),
		cf:     config,
		ch:     make(chan logger.IData, MsgChanLen),
		msgBuf: make([]logger.IData, 0, MaxMsgBufLen),
		buf:    &bytes.Buffer{},
	}

	e.w = e.getWriter()

	go func() {
		defer func() {
			if rev := recover(); rev != nil {
				fmt.Printf("日志写盘异常%v", rev)
			}
		}()
		e.flush()
	}()

	return e
}

func (e *entity) base(level logger.Level) logger.IData {
	return logger.InitData().
		SetLevel(level).
		Common(logger.LevelTag, base.LevelInfo(level)).
		Common(logger.CallerTag, base.GetCallerInfo(e.GetSkip())).
		Common(logger.TimeTag, time.Now().Format("2006-01-02T15:04:05.000Z0700"))
}

func (e *entity) Debugf(msg string, values ...interface{}) {
	e.output(e.base(logger.Debug).Common(logger.MsgTag, fmt.Sprintf(msg, values...)))
}

func (e *entity) Infof(msg string, values ...interface{}) {
	e.output(e.base(logger.Info).Common(logger.MsgTag, fmt.Sprintf(msg, values...)))
}

func (e *entity) Warnf(msg string, values ...interface{}) {
	e.output(e.base(logger.Warn).Common(logger.MsgTag, fmt.Sprintf(msg, values...)))
}

func (e *entity) Errorf(msg string, values ...interface{}) {
	e.output(e.base(logger.Error).Common(logger.MsgTag, fmt.Sprintf(msg, values...)))
}

func (e *entity) Debug(msg string, values ...logger.Field) {
	fields := logger.NewFields()
	for _, value := range values {
		value(fields)
	}
	e.output(e.base(logger.Debug).Common(logger.MsgTag, msg).Map(fields.Data()))
}

func (e *entity) Info(msg string, values ...logger.Field) {
	fields := logger.NewFields()
	for _, value := range values {
		value(fields)
	}
	e.output(e.base(logger.Info).Common(logger.MsgTag, msg).Map(fields.Data()))
}

func (e *entity) Warn(msg string, values ...logger.Field) {
	fields := logger.NewFields()
	for _, value := range values {
		value(fields)
	}
	e.output(e.base(logger.Warn).Common(logger.MsgTag, msg).Map(fields.Data()))
}

func (e *entity) Error(msg string, values ...logger.Field) {
	fields := logger.NewFields()
	for _, value := range values {
		value(fields)
	}
	e.output(e.base(logger.Error).Common(logger.MsgTag, msg).Map(fields.Data()))
}

// 监听刷盘情况
func (e *entity) flush() {
	// 监听信号
	c := make(chan os.Signal)
	// 监听信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 注册定时器
	ticker := time.NewTicker(FlushTimeInterval)

	for {
		select {
		case mes := <-e.ch:
			e.msgBuf = append(e.msgBuf, mes)
			if len(e.msgBuf) == MaxMsgBufLen {
				if err := e.brush(); err != nil {
					fmt.Printf("缓冲区上限刷盘，错误:%v\n", err)
				}
			}
		case <-ticker.C:
			if err := e.brush(); err != nil {
				fmt.Printf("缓冲区定时刷盘，错误:%v\n", err)
			}
		case <-c:
			if err := e.brush(); err != nil {
				fmt.Printf("退出刷盘，错误:%v\n", err)
			}
		}
	}
}

// 批量将buf内容刷新到日志文件
func (e *entity) brush() (err error) {
	e.buf.Reset()
	for _, mes := range e.msgBuf {
		e.buf.Write(mes.Marshal())
		e.buf.WriteByte('\n')
	}

	// 重置
	e.msgBuf = make([]logger.IData, 0, MaxMsgBufLen)

	// 检查
	if e.buf.Len() == 0 {
		return
	}

	// 写入日志文件
	_, err = io.Copy(e.w, e.buf)
	if err != nil {
		fmt.Println("写入日志文件失败,", err)
		return
	}

	return
}

func (e *entity) output(data logger.IData) {
	if data.GetLevel() > base.Level(e.cf.GetConfig().Level) {
		e.ch <- data
	}
}

func (e *entity) getWriter() io.Writer {
	output := strings.Split(e.cf.GetConfig().Out, ",")
	if len(output) == 0 {
		output = append(output, "std")
	}

	otw := make([]io.Writer, 0)
	for _, ot := range output {

		fc := e.cf.GetFileConfig()

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
