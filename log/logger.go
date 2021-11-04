package log

import (
	"fmt"

	"github.com/antsgo/antsgo/conf"
	"github.com/sirupsen/logrus"
)

func NewLogger(conf conf.Conf) map[string]*logrus.Logger {
	logger := make(map[string]*logrus.Logger)
	for name, value := range conf.Logger {
		switch value.Type {
		case "file":
			logger[name] = newLoggerFile(value)
			break
		case "kafka":
			fmt.Println("kafka")
			break
		}
	}
	return logger
}
