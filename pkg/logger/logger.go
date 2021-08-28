package logger

type Logger interface {
	Fatal(msg string, keyVal ...KeyVal)
	Error(msg string, keyVal ...KeyVal)
	Debug(msg string, keyVal ...KeyVal)
	Info(msg string, keyVal ...KeyVal)
}

type LoggerFactory interface {
	NewLogger() Logger
}

func NewLogger(factory LoggerFactory) Logger {
	return factory.NewLogger()
}
