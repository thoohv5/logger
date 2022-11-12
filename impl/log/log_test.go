package log

import (
	"sync"
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
	l.Debug("22222", logger.FieldInt32("a", 2))
	l.Info("3333", logger.FieldInt32("a", 2))
}

func Test_entity_Debug_go(t *testing.T) {
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
	wg := sync.WaitGroup{}
	for i := 0; i < 100000000; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			l.Debug("22222", logger.FieldInt32("a", 2))
			l.Debugf("22222:%v", 111)
		}()
		go func() {
			defer wg.Done()
			l.Info("3333", logger.FieldInt32("a", 2))
			l.Infof("3333:%v", 2222)
		}()
	}
	wg.Wait()
}

func BenchmarkRepeat(b *testing.B) {
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

	l.Info("xxxxx", logger.FieldInt32("k", 1))
}
