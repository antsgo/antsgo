package log

import (
	"github.com/antsgo/antsgo/conf"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func newLoggerFile(conf conf.ConfLogger) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(&lumberjack.Logger{
		Filename:   conf.Path,
		MaxSize:    conf.MaxSize, // 文件大小[单位mb]
		MaxBackups: conf.FileNum, // 保留文件个数
		MaxAge:     conf.DayNum,  // 保留天数
		Compress:   true,         // 是否压缩日志
	})
	return logger
}
