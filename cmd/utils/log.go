package utils

import (
	"github.com/offerm/lnagent/logger"
	"github.com/sirupsen/logrus"
)

func init() {
	config := &logger.Config{
		AppName:      "lnagent",
		FileName:     "/tmp/lnagent-logs/lnagent.log",
		FileLevel:    logrus.TraceLevel,
		ConsoleLevel: logrus.InfoLevel,
		MaxSize:      100,
		MaxBackups:   3,
		MaxAge:       30,
	}
	logger.Init(config, "v0.0.1")
}
