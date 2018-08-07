package worker

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/DeslumTeam/shkaff/apps/statsender"
	"github.com/DeslumTeam/shkaff/drivers/mongodb"
	"github.com/DeslumTeam/shkaff/internal/databases"
	"github.com/DeslumTeam/shkaff/internal/logger"
	"github.com/DeslumTeam/shkaff/internal/structs"

	"github.com/op/go-logging"
)

const (
	WORKER_SLEEP_TIMEOUT = 100 * time.Millisecond
)

type worker struct {
	workerWG     sync.WaitGroup
	dumpChan     chan *structs.Task
	databaseName string
	log          *logging.Logger
	stat         *statsender.StatSender
	dbDriver     databases.DatabaseDriver
}

func (w *worker) Run() {
	w.stat = statsender.Run()
	w.dumpChan = make(chan *structs.Task, 1)
	w.log = logger.GetLogs("Worker")
	w.workerWG.Add(1)
	go func() {
		defer w.workerWG.Done()
		w.proc()
	}()
}

func (w *worker) Send(task *structs.Task) (err error) {
	w.databaseName = task.DBType
	w.dbDriver, err = w.getDatabaseType()
	if err != nil {
		w.log.Error(err)
		return
	}

	w.dumpChan <- task
	return
}

func (w *worker) Stop() {
	close(w.dumpChan)
	w.workerWG.Done()
}

func (w *worker) proc() {
	for {
		task, ok := <-w.dumpChan
		if !ok {
			break
		}

		w.stat.SendStatMessage(structs.StartDump, task.UserID, task.DBID, task.TaskID, nil)
		err := w.dbDriver.Dump(task)
		if err != nil {
			w.stat.SendStatMessage(structs.FailDump, task.UserID, task.DBID, task.TaskID, err)
			w.log.Error(err)
			continue
		}

		w.stat.SendStatMessage(structs.SuccessDump, task.UserID, task.DBID, task.TaskID, nil)
		w.stat.SendStatMessage(structs.NewRestore, task.UserID, task.DBID, task.TaskID, nil)
		err = w.dbDriver.Restore(task)
		if err != nil {
			w.stat.SendStatMessage(structs.FailRestore, task.UserID, task.DBID, task.TaskID, err)
			w.log.Error(err)
			continue
		}

		w.stat.SendStatMessage(structs.SuccessRestore, task.UserID, task.DBID, task.TaskID, err)
		time.Sleep(WORKER_SLEEP_TIMEOUT)
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
