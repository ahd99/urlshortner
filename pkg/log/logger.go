package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	test11()
}

func test11() {
	errToConsoleLogger := getErrorToConsoleLogger()
	tee := zapcore.NewTee(errToConsoleLogger)

	logger := zap.New(tee, zap.AddCaller(), zap.Development())
	logger.Error("ERRRRRRor")
}

func getErrorToConsoleLogger() zapcore.Core {
	//levelEnabler := zap.DebugLevel
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l > zap.WarnLevel
	})

	wsync := zapcore.AddSync(os.Stdout)

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


