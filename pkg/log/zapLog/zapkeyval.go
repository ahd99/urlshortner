package zapLog

import (
	"go.uber.org/zap/zapcore"
)

type ZapKeyVal struct {
	zapcore.Field
}

func (kv *ZapKeyVal) GetKey() string {
	return kv.Key
}

func (kv *ZapKeyVal) GetValue() interface{} {
	return nil
}
