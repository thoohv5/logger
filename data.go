package logger

import (
	"encoding/json"
	"sync"
)

type (
	IData interface {
		Append(key string, values ...interface{}) IData
		Marshal() string
	}
	Data struct {
		data sync.Map
	}
)

func InitData() *Data {
	return &Data{}
}

func (ld *Data) Append(key string, values ...interface{}) IData {
	if len(values) == 1 {
		ld.data.Store(key, values[0])
	} else {
		ld.data.Store(key, values)
	}
	return ld
}

func (ld *Data) Marshal() string {
	m := make(map[string]interface{})
	ld.data.Range(func(key, value interface{}) bool {
		m[key.(string)] = value
		return true
	})
	bs, _ := json.Marshal(m)
	return string(bs)
}
