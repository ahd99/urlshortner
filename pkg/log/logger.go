package logger

type Logger interface {
	Error(msg string, keyVal ...KeyVal)
	Debug(msg string, keyVal ...KeyVal)
	Info(msg string, keyVal ...KeyVal)
}

type LoggerFactory interface {
	NewLogger() *Logger
}

var loggerFactory LoggerFactory

func SetLoggerFactory(f LoggerFactory) {
	loggerFactory = f
}

func NewLogger() *Logger {
	return loggerFactory.NewLogger()
}
