package logger

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type ILogger interface {
	Debug(msg string, values ...interface{})
	Info(msg string, values ...interface{})
	Warn(msg string, values ...interface{})
	Error(msg string, values ...interface{})
}
