package logger

import "github.com/thoohv5/logger/util"

const (
	_defaultPath       = "./"
	_defaultFileName   = "default.log"
	_defaultMaxSize    = 128
	_defaultMaxBackups = 30
	_defaultMaxAge     = 7
	_defaultCompress   = true
)

// Config 日志配置
type Config struct {
	// 日志: std, file 多种类型按照逗号隔开
	Out string `toml:"out"`
	// 日志类别: debug, warn, info，error
	Level string `toml:"level"`
	// Out包含file，需要配置
	File *File `toml:"file"`
}

// File 日志文件类别配置
type File struct {
	// 日志文件目录 项目目录
	Path string `toml:"path"`
	// 日志文件默认名称 default.log
	FileName string `toml:"file_name"`
	// 最大文件容量（单位:M），默认：128
	MaxSize int `toml:"max_size"`
	// 最大的备份数量，默认：30
	MaxBackups int `toml:"max_backups"`
	// 最大备份天数，默认：0，没有限制
	MaxAge int `toml:"max_age"`
	// 是否gzip压缩
	Compress bool `toml:"compress"`
}

func (c *Config) GetFileConfig() *File {
	if c.File == nil {
		c.File = &File{
			Path:       util.AbPath(_defaultPath),
			FileName:   _defaultFileName,
			MaxSize:    _defaultMaxSize,
			MaxBackups: _defaultMaxBackups,
			MaxAge:     _defaultMaxAge,
			Compress:   _defaultCompress,
		}
	}
	if c.File.Path == "" {
		c.File.Path = util.AbPath(_defaultPath)
	}
	if c.File.FileName == "" {
		c.File.FileName = _defaultFileName
	}
	return c.File
}

func (c *Config) GetConfig() *Config {
	if c.Out == "" {
		c.Out = "std"
	}
	if c.Level == "" {
		c.Level = "info"
	}
	return c
}
