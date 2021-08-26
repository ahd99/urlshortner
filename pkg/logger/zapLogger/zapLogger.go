package zapLogger

import (
	"time"

	"github.com/ahd99/urlshortner/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"os"
)

type ZapKeyValFactory struct{}

func (f *ZapKeyValFactory) Int64(key string, val int64) logger.KeyVal {
	return zap.Int64(key, val)
}

func (f *ZapKeyValFactory) Int(key string, val int) logger.KeyVal {
	return zap.Int(key, val)
}

func (f *ZapKeyValFactory) String(key string, val string) logger.KeyVal {
	return zap.String(key, val)
}

type ZapLoggerFactory struct{}

func (f ZapLoggerFactory) NewLogger() logger.Logger {
	logger1 := newLogger()
	return logger1
}

type zapLogger struct {
	zaplogger *zap.Logger
}

func (l *zapLogger) Error(msg string, keyVal ...logger.KeyVal) {
	zapfields := make([]zap.Field, len(keyVal))
	for i := 0; i < len(keyVal); i++ {
		zapfields[i] = keyVal[i].(zap.Field)
	}
	l.zaplogger.Error(msg, zapfields...)
}

func (l *zapLogger) Debug(msg string, keyVal ...logger.KeyVal) {
	zapfields := make([]zap.Field, len(keyVal))
	for i := 0; i < len(keyVal); i++ {
		zapfields[i] = keyVal[i].(zap.Field)
	}
	l.zaplogger.Debug(msg, zapfields...)
}

func (l *zapLogger) Info(msg string, keyVal ...logger.KeyVal) {
	zapfields := make([]zap.Field, len(keyVal))
	for i := 0; i < len(keyVal); i++ {
		zapfields[i] = keyVal[i].(zap.Field)
	}
	l.zaplogger.Info(msg, zapfields...)
}

func stdLogger(pretty bool) zapcore.Core {
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.DebugLevel
	})
	wsync := zapcore.Lock(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("2006-01-02 15:04:05 Z0700"))
	}

	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	var encoder zapcore.Encoder
	if pretty {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, wsync, levelEnabler)

	return core
}

func newLogger() logger.Logger {
	stdLogger := stdLogger(true)
	tee := zapcore.NewTee(stdLogger)
	logger1 := zap.New(tee, zap.AddCaller(), zap.Development())
	zaplogger := zapLogger{logger1}
	return &zaplogger
}
