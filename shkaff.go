package main

import (
	"log"
	"os"
	"shkaff/apps/api"
	"shkaff/apps/operator"
	"shkaff/apps/statsender"
	"shkaff/apps/worker"
	"shkaff/internal/fork"
	"shkaff/internal/options"
)

type Creater interface {
	Init(action string, cfg options.ShkaffConfig) *Service
}
type Service interface {
	Run()
	Stop()
}

type shkaff struct{}

func (self *shkaff) Init(action string) (srv Service) {
	switch action {
	case "Operator":
		srv = operator.InitOperator()
	case "Worker":
		srv = worker.InitWorker()
	case "StatWorker":
		srv = statsender.InitStatSender()
	case "API":
		srv = api.InitAPI()
	default:
		log.Fatalf("Unknown Shkaff service name %s\n", action)
	}
	return
}

func startShkaff() {
	servicesName := []string{"Operator", "Worker", "StatWorker", "API"}
	shkf := new(shkaff)
	for _, name := range servicesName {
		s := shkf.Init(name)
		go s.Run()
	}
}

func main() {
	daemon, err := fork.InitDaemon()
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
	status, err := daemon.Run(startShkaff)
	if err != nil {
		log.Println(status, "\nError: ", err)
		os.Exit(3)
	}
	log.Println(status)
}
