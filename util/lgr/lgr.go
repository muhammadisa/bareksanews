package lgr

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
	"time"
)

type LogType string

const (
	LogInfo LogType = `logInfo`
	LogWarn LogType = `logWarn`
	LogErr  LogType = `logErr`
	LogData LogType = `logData`
)

func Create(serviceName string) log.Logger {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	//logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(
		logger,
		"service", serviceName,
		"time", log.DefaultTimestampUTC,
		"caller", log.Caller(3),
	)
	currentTime := time.Now()
	level.Info(logger).Log(LogInfo, fmt.Sprintf("the service is started at %s", currentTime))
	return logger
}
