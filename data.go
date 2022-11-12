package logger

import (
	"bytes"
	"fmt"
	"sort"
)

type (
	IData interface {
		Common(key string, values ...interface{}) IData
		Custom(key string, values ...interface{}) IData
		Map(value map[string]interface{}) IData

		Marshal() []byte
		SetLevel(level Level) IData
		GetLevel() Level
	}
	data struct {
		level Level
		keys  []string
		data  map[string]interface{}
	}
)

func InitData() IData {
	return &data{
		keys: make([]string, 0),
		data: make(map[string]interface{}),
	}
}

func (ld *data) Common(key string, values ...interface{}) IData {
	if len(values) == 1 {
		ld.data[key] = values[0]
	} else {
		ld.data[key] = values
	}
	return ld
}

func (ld *data) Custom(key string, values ...interface{}) IData {
	ld.keys = append(ld.keys, key)
	if len(values) == 1 {
		ld.data[key] = values[0]
	} else {
		ld.data[key] = values
	}
	return ld
}

func (ld *data) Map(items map[string]interface{}) IData {
	for k, v := range items {
		ld.keys = append(ld.keys, k)
		ld.data[k] = v
	}
	return ld
}

func (ld *data) Marshal() []byte {
	// bs, _ := json.Marshal(ld.data)

	sort.Strings(ld.keys)

	keys := make([]string, 0, len(ld.data))

	keys = append([]string{
		CallerTag,
		LevelTag,
		TimeTag,
		MsgTag,
	}, ld.keys...)

	buf := &bytes.Buffer{}
	buf.WriteByte('{')
	for idx, key := range keys {
		if idx > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(fmt.Sprintf("\"%s\":\"%v\"", key, ld.data[key]))
	}
	buf.WriteByte('}')

	return buf.Bytes()
}

func (ld *data) SetLevel(level Level) IData {
	ld.level = level
	return ld
}

func (ld *data) GetLevel() Level {
	return ld.level
}
