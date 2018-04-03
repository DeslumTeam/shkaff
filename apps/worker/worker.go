package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/DeslumTeam/shkaff/apps/statsender"
	"github.com/DeslumTeam/shkaff/drivers/mongodb"
	"github.com/DeslumTeam/shkaff/drivers/rmq/consumer"
	"github.com/DeslumTeam/shkaff/internal/databases"
	"github.com/DeslumTeam/shkaff/internal/logger"
	"github.com/DeslumTeam/shkaff/internal/structs"

	logging "github.com/op/go-logging"
)

type workersStarter struct {
	workers    map[string]map[int]*worker // map[driverName]map[taskID]*Worker
	workRabbit *consumer.RMQ
	log        *logging.Logger
}

type worker struct {
	workerWG     sync.WaitGroup
	dumpChan     chan *structs.Task
	databaseName string
	log          *logging.Logger
	stat         *statsender.StatSender
	dbDriver     databases.DatabaseDriver
}

func InitWorker() (ws *workersStarter) {
	return &workersStarter{
		workRabbit: consumer.InitAMQPConsumer(),
		workers:    make(map[string]map[int]*worker),
		log:        logger.GetLogs("WorkerManager"),
	}
}

func (ws *workersStarter) Run() {
	var task *structs.Task
	ws.log.Infof("Run WorkersManager")
	for message := range ws.workRabbit.Msgs {
		err := json.Unmarshal(message.Body, &task)
		if err != nil {
			ws.log.Error(err)
			continue
		}
		if ws.workers[task.DBType] == nil {
			ws.workers[task.DBType] = make(map[int]*worker)
		}
		if ws.workers[task.DBType][task.TaskID] == nil {
			ws.workers[task.DBType][task.TaskID] = new(worker)
			ws.workers[task.DBType][task.TaskID].Run()
		}
		err = ws.workers[task.DBType][task.TaskID].Send(task)
		if err != nil {
			ws.log.Error(err)
			continue
		}
	}
}

func (ws *workersStarter) Stop() {
	ws.workRabbit.Channel.Close()
	return
}

func (w *worker) Run() {
	w.stat = statsender.Run()
	w.dumpChan = make(chan *structs.Task)
	w.workerWG.Add(1)
	go w.proc()
	w.workerWG.Wait()
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
		w.stat.SendStatMessage(3, task.UserID, task.DBID, task.TaskID, nil)
		err := w.dbDriver.Dump(task)
		if err != nil {
			w.stat.SendStatMessage(5, task.UserID, task.DBID, task.TaskID, err)
			w.log.Error(err)
			return
		}
		w.stat.SendStatMessage(4, task.UserID, task.DBID, task.TaskID, nil)
		w.stat.SendStatMessage(6, task.UserID, task.DBID, task.TaskID, nil)
		err = w.dbDriver.Restore(task)
		if err != nil {
			w.stat.SendStatMessage(8, task.UserID, task.DBID, task.TaskID, err)
			w.log.Error(err)
			return
		}
		w.stat.SendStatMessage(7, task.UserID, task.DBID, task.TaskID, err)
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
