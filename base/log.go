package base

import (
	"github.com/hhq163/logger"
)

var Log logger.Logger

func LogInit(debug bool, svcName string) {
	Log = logger.NewDefaultLogger()
	if debug { //非json格式
		config := logger.NewDevelopmentConfig()
		logger.SetConfig(config)
	}
}
