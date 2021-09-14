package main

import (
	"github.com/ahd99/urlshortner/internal/monitoring/client"
	"github.com/ahd99/urlshortner/pkg/logger"
	"github.com/ahd99/urlshortner/pkg/logger/zapLogger"
)

var logger1 logger.Logger

func main() {
	logger1 = initLogger()

	client.Init("185.235.40.218:8091", logger1)
	defer client.Cleanup()

}

func initLogger() logger.Logger {
	logger.SetKeyValuePairFactory(&zapLogger.ZapKeyValFactory{})
	logger1 := logger.NewLogger(zapLogger.ZapLoggerFactory{})
	return logger1
}
