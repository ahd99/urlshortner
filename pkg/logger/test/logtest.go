package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	//test1()
	//test2()
}

// simplest use
func test0() {
	logger := zap.NewExample()
	//logger := zap.NewProduction()
	//logger := zap.NewDevelopment()
	defer logger.Sync()

	logger.Debug("THis is ,essage",
		zap.String("strField", "ali"),
		zap.Int("weight", 3),
		zap.Duration("time", 10*time.Second))

	sugar := logger.Sugar()
	sugar.Infow("log message", "time", 3, "weight", 10) // {"level":"info","msg":"log message","time":3,"weight":10}
	sugar.Infof("log message %d %d", 3, 10)             // {"level":"info","msg":"log message 3 10"}
	sugar.Info("time: ", 3, "  weight: ", 10)           //{"level":"info","msg":"time: 3  weight: 10"}

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
