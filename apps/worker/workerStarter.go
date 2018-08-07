package worker

import (
	"encoding/json"
	"github.com/op/go-logging"
	"github.com/DeslumTeam/shkaff/internal/logger"
	"github.com/DeslumTeam/shkaff/internal/structs"
	"github.com/DeslumTeam/shkaff/drivers/rmq/consumer"
)

type workersStarter struct {
	workers    map[string]map[int]*worker // map[driverName]map[taskID]*Worker
	workRabbit *consumer.RMQ
	log        *logging.Logger
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
	ws.workRabbit.InitConnection("mongodb")
	ws.log.Info("Start WorkersManager")
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

		message.Ack(true)
	}
}

func (ws *workersStarter) Stop() {
	ws.workRabbit.Channel.Close()
	return
}
