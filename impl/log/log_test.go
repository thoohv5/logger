package log

import (
	"testing"

	"github.com/thoohv5/logger"
	"github.com/thoohv5/logger/util"
)

func Test_entity_Debug(t *testing.T) {
	l := New(&logger.Config{
		Out:   "std,file",
		Level: "debug",
		File: &logger.File{
			Path:       util.AbPath("../../log"),
			FileName:   "default.log",
			MaxSize:    0,
			MaxBackups: 0,
			MaxAge:     0,
			Compress:   false,
		},
	})
	l.Debug("22222", map[string]interface{}{"a": 2})
	l.Info("3333", map[string]interface{}{"a": 2})
}
