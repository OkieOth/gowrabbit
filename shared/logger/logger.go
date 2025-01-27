package logger

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}
