package logger

import (
	"os"

	"github.com/op/go-logging"
)

func GetLogs(appName string) (log *logging.Logger) {
	log = logging.MustGetLogger(appName)
	var format = logging.MustStringFormatter(
		`%{color} %{time:15:04:05} %{module:-12s} %{shortfunc:-9s} %{level:-5s}%{color:reset} %{message}`,
	)
	logHandler := logging.NewLogBackend(os.Stdout, appName, 0)
	logFormatter := logging.NewBackendFormatter(logHandler, format)
	logging.SetBackend(logFormatter)
	return log
}
