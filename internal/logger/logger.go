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
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backend2Formatter)
	return
}
