package logger

import (
	"os"

	"github.com/op/go-logging"
)

const(
	LOG_FORMAT = `%{color}%{level:-5s} %{time:15:04:05} %{module:-15s} %{shortfunc:-9s} %{color:reset} %{message}`
)

func GetLogs(appName string) (log *logging.Logger) {
	log = logging.MustGetLogger(appName)
	logHandler := logging.NewLogBackend(os.Stdout, "", 0)
	logFormatter := logging.NewBackendFormatter(logHandler, logging.MustStringFormatter(LOG_FORMAT))
	logging.SetBackend(logFormatter)
	return log
}
