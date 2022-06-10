package logger

import "encoding/json"

type (
	IData interface {
		Append(key string, values ...interface{}) IData
		Marshal() string
	}
	Data struct {
		data map[string]interface{}
	}
)

func InitData() *Data {
	return &Data{
		data: make(map[string]interface{}),
	}
}

func (ld *Data) Append(key string, values ...interface{}) IData {
	if len(values) == 1 {
		ld.data[key] = values[0]
	} else {
		ld.data[key] = values
	}
	return ld
}

func (ld *Data) Marshal() string {
	bs, _ := json.Marshal(ld.data)
	return string(bs)
}
