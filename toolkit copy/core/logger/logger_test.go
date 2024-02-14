package logger_test

import (
	"testing"

	"gitlab.com/tuneverse/toolkit/core/logger"
)

func TestLogger(t *testing.T) {
	t.Run("init-logger", func(t *testing.T) {
		logger.InitLogger(&logger.ClientOptions{
			Service:             "service",
			LogLevel:            "info",
			IncludeRequestDump:  true,
			IncludeResponseDump: true,
			JSONFormater:        true,
		}, &logger.FileMode{
			LogfileName: "error.log",
		})

		logger.Log().Error("demo %v", "1")
		logger.Log().Error("demo-001 %v", "1")
	})
}
