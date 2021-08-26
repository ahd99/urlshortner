package logger

type KeyVal interface{}

type KeyValFactory interface {
	Int64(key string, val int64) KeyVal
	Int(key string, val int) KeyVal
	String(key string, val string) KeyVal
}

var kvFactory KeyValFactory

func SetKeyValuePairFactory(f KeyValFactory) {
	kvFactory = f
}

func Int64(key string, val int64) KeyVal {
	return kvFactory.Int64(key, val)
}

func Int(key string, val int) KeyVal {
	return kvFactory.Int(key, val)
}

func String(key string, val string) KeyVal {
	return kvFactory.String(key, val)
}
