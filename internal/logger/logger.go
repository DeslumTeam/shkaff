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
	// file, err := os.OpenFile("shkaff.log",
	// 	os.O_CREATE|os.O_WRONLY|os.O_APPEND,
	// 	0666)
	// if err != nil {
	// 	fmt.Printf("Failed to open log file. Error: %s\n", err)
	// 	os.Exit(-1)
	// }
	logHandler := logging.NewLogBackend(os.Stdout, appName, 0)
	logFormatter := logging.NewBackendFormatter(logHandler, format)
	logging.SetBackend(logFormatter)
	return log
}
