package base

import (
	"github.com/hhq163/logger"
)

var Log logger.Logger

func LogInit(debug bool, svcName string) {
	if debug {
		cfg := logger.NewDevelopmentConfig()
		// cfg.Encoding = "json"
		cfg.OutputPaths = append(cfg.OutputPaths, "access_log.txt")
		Log = logger.NewMyLogger(cfg)

	} else {
		cfg := logger.NewProductionConfig()
		Log = logger.NewMyLogger(cfg)
	}
}
