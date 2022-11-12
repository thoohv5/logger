package logger

type (
	IFields interface {
		Set(key string, val interface{}) IFields
		Data() map[string]interface{}
	}
	fields struct {
		data map[string]interface{}
	}
	Field func(IFields)
)

func NewFields() IFields {
	return &fields{
		data: make(map[string]interface{}),
	}
}

func (fs *fields) Data() map[string]interface{} {
	return fs.data
}
func (fs *fields) Set(key string, val interface{}) IFields {
	fs.data[key] = val
	return fs
}

func FieldInt32(key string, val int32) Field {
	return func(fs IFields) {
		fs.Set(key, val)
	}
}

func FieldInt64(key string, val int64) Field {
	return func(fs IFields) {
		fs.Set(key, val)
	}
}

func FieldString(key string, val string) Field {
	return func(fs IFields) {
		fs.Set(key, val)
	}
}

func FieldMap(items map[string]interface{}) Field {
	return func(fs IFields) {
		for k, v := range items {
			fs.Set(k, v)
		}
	}
}
