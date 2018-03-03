package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"shkaff/apps/statsender"
	"shkaff/drivers/maindb"
	"shkaff/drivers/mongodb"
	"shkaff/drivers/rmq/consumer"
	"shkaff/internal/databases"
	"shkaff/internal/logger"
	"shkaff/internal/options"
	"shkaff/internal/structs"
	"sync"

	logging "github.com/op/go-logging"
)

type workersStarter struct {
	workerWG    sync.WaitGroup
	workers     []*worker
	workerCount int
}

type worker struct {
	dumpChan     chan string
	databaseName string
	postgres     *maindb.PSQL
	workRabbit   *consumer.RMQ
	stat         *statsender.StatSender
	log          *logging.Logger
}

func InitWorker() (ws *workersStarter) {
	cfg := options.InitControlConfig()
	ws = new(workersStarter)
	ws.workerCount = cfg.WORKERS["mongodb"]
	stat := statsender.Run()
	ws.workerWG = sync.WaitGroup{}
	for i := 0; i < ws.workerCount; i++ {
		worker := &worker{
			databaseName: "mongodb",
			dumpChan:     make(chan string, 100),
			postgres:     maindb.InitPSQL(),
			workRabbit:   consumer.InitAMQPConsumer(),
			stat:         stat,
			log:          logger.GetLogs("Worker"),
		}
		ws.workers = append(ws.workers, worker)
	}
	return
}

func (ws *workersStarter) Run() {
	ws.workerWG.Add(ws.workerCount)
	for _, w := range ws.workers {
		go w.worker()
	}
	ws.workerWG.Wait()
}

func (ws *workersStarter) Stop() {
	for w := 0; w < ws.workerCount; w++ {
		ws.workerWG.Done()
	}
}

func (w *worker) worker() {
	var task *structs.Task
	w.workRabbit.InitConnection(w.databaseName)
	dbDriver, err := w.getDatabaseType()
	if err != nil {
		w.log.Error(err)
		return
	}
	w.log.Infof("Start Worker for %s\n", w.databaseName)
	for message := range w.workRabbit.Msgs {
		err := json.Unmarshal(message.Body, &task)
		if err != nil {
			w.log.Error(err)
			continue
		}
		w.stat.SendStatMessage(3, task.UserID, task.DBID, task.TaskID, nil)
		err = dbDriver.Dump(task)
		if err != nil {
			w.stat.SendStatMessage(5, task.UserID, task.DBID, task.TaskID, err)
			w.log.Error(err)
			continue
		}
		w.stat.SendStatMessage(4, task.UserID, task.DBID, task.TaskID, nil)
		w.stat.SendStatMessage(6, task.UserID, task.DBID, task.TaskID, nil)
		err = dbDriver.Restore(task)
		if err != nil {
			w.stat.SendStatMessage(8, task.UserID, task.DBID, task.TaskID, err)
			w.log.Error(err)
			continue
		}
		w.stat.SendStatMessage(7, task.UserID, task.DBID, task.TaskID, err)
		message.Ack(false)
	}
}

func (w *worker) getDatabaseType() (dbDriver databases.DatabaseDriver, err error) {
	switch w.databaseName {
	case "mongodb":
		dbDriver = mongodb.InitDriver()
		return dbDriver, nil
	default:
		answer := fmt.Sprintf("Driver %s not found", w.databaseName)
		return nil, errors.New(answer)
	}
}
