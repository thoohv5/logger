package logger

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type ILogger interface {
	Debugf(msg string, values ...interface{})
	Infof(msg string, values ...interface{})
	Warnf(msg string, values ...interface{})
	Errorf(msg string, values ...interface{})

	Debug(msg string, values ...Field)
	Info(msg string, values ...Field)
	Warn(msg string, values ...Field)
	Error(msg string, values ...Field)
}
