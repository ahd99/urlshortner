package logger

type KeyVal interface {
	GetKey() string
	GetValue() interface{}
}

type KeyValFactory interface {
	Int64(key string, val int64) KeyVal
	Int(key string, val int) KeyVal
	String(key string, val string) KeyVal
}

var factory KeyValFactory

func SetKeyValuePairFactory(f KeyValFactory) {
	factory = f
}

func Int64(key string, val int64) KeyVal {
	return factory.Int64(key, val)
}

func Int(key string, val int) KeyVal {
	return factory.Int(key, val)
}

func IntString(key string, val string) KeyVal {
	return factory.String(key, val)
}
