package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	//test1()
	test2()
}

func test1() {
	errToConsoleLogger := getErrorToConsoleLogger()
	tee := zapcore.NewTee(errToConsoleLogger)

	logger := zap.New(tee, zap.AddCaller(), zap.Development())
	defer logger.Sync()
	logger.Error("ERRRRRRor")

	zap.ReplaceGlobals(logger)
	zap.L().Error("Global Errrrrror")
	zap.S().Error("Global sugared Errrrrror") // zap.S() returns sugar logger
}

func getErrorToConsoleLogger() zapcore.Core {
	//levelEnabler := zap.DebugLevel
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l > zap.WarnLevel
	})

	wsync := zapcore.AddSync(os.Stdout)
	//wsync := os.Stdout	// because os.Stdout implementes zaocore.WriteSyncer (its implement sync() method) so we can use it without zapcore.AddSync() method
	wsync = zapcore.Lock(wsync) // zap official doc for Lock(): Lock wraps a WriteSyncer in a mutex to make it safe for concurrent use. In particular, *os.Files must be locked before use.

	encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("2006-01-02 15:04:05 Z0700"))
	}

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, wsync, levelEnabler)

	return core
}
