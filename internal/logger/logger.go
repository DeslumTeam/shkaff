package logger

import (
	"os"

	"github.com/op/go-logging"
)

func GetLogs(appName string) (log *logging.Logger) {
	log = logging.MustGetLogger(appName)
	var format = logging.MustStringFormatter(
		`%{color}%{level:-5s} %{time:15:04:05} %{module:-15s} %{shortfunc:-9s} %{color:reset} %{message}`,
	)
	logHandler := logging.NewLogBackend(os.Stdout, "", 0)
	logFormatter := logging.NewBackendFormatter(logHandler, format)
	logging.SetBackend(logFormatter)
	return log
}
