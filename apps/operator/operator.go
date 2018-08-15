package operator

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/DeslumTeam/shkaff/apps/statsender"
	"github.com/DeslumTeam/shkaff/drivers/cache"
	"github.com/DeslumTeam/shkaff/drivers/maindb"
	"github.com/DeslumTeam/shkaff/drivers/mongodb"
	"github.com/DeslumTeam/shkaff/drivers/rmq/producer"
	"github.com/DeslumTeam/shkaff/internal/consts"
	"github.com/DeslumTeam/shkaff/internal/logger"
	"github.com/DeslumTeam/shkaff/internal/structs"

	_ "github.com/lib/pq"
	"github.com/op/go-logging"
)

const (
	DATATIME_FORMAT         = "2006-01-02 15:03"
	AGGREGATOR_ITER_TIMEOUT = 10 * time.Second
	MONGO                   = "mongodb"
)

type Operator struct {
	tasksChan  chan structs.Task
	operatorWG sync.WaitGroup
	postgres   *maindb.PSQL
	rabbit     *producer.RMQ
	taskCache  *cache.Cache
	log        *logging.Logger
	stat       *statsender.StatSender
}

func InitOperator() (o *Operator) {
	log := logger.GetLogs("Operator")
	producer, err := producer.InitAMQPProducer("mongodb")
	if err != nil {
		log.Fatalf("Operator error %v", err)
	}
	o = &Operator{
		postgres:   maindb.InitPSQL(),
		rabbit:     producer,
		tasksChan:  make(chan structs.Task),
		log:        log,
		operatorWG: sync.WaitGroup{},
		stat:       statsender.Run(),
	}

	o.taskCache, err = cache.InitCacheDB()
	if err != nil {
		o.log.Fatal(err)
	}

	return
}

// Run Start Operator service
func (o *Operator) Run() {
	o.operatorWG.Add(1)
	go func() {
		defer o.operatorWG.Done()
		o.log.Info("Start Aggregator")
		o.aggregator()
	}()

	o.operatorWG.Add(1)
	go func() {
		defer o.operatorWG.Done()
		o.log.Info("Start TaskSender")
		o.taskSender()
	}()

	o.operatorWG.Wait()
}

// Stop Operator service
func (o *Operator) Stop() {
	for i := 0; i < 2; i++ {
		o.operatorWG.Done()
	}
	o.log.Info("Stop Operator")
}

func (o *Operator) taskSender() {
	for {
		task, ok := <-o.tasksChan
		if !ok {
			break
		}

		messages, err := o.getMessagesByDatabaseType(task)
		if err != nil {
			o.log.Errorf("getMessagesByDatabaseType %v \n Error %v", task, err)
			continue
		}

		for _, msg := range messages {
			o.log.Infof("-->> %v", msg)
			o.stat.SendStatMessage(structs.NewOperator, task.UserID, task.DBID, task.TaskID, nil)
			body, err := json.Marshal(msg)
			if err != nil {
				o.stat.SendStatMessage(structs.FailOperator, task.UserID, task.DBID, task.TaskID, err)
				o.log.Error(err)
				continue
			}
			o.log.Infof("->> %v", body)
			err = o.rabbit.Publish(body)
			if err != nil {
				o.stat.SendStatMessage(structs.FailOperator, task.UserID, task.DBID, task.TaskID, err)
				o.log.Error(err)
				continue
			}

			o.stat.SendStatMessage(structs.StartOperator, task.UserID, task.DBID, task.TaskID, nil)
		}
	}
}

func (o *Operator) aggregator() {
	psqlUpdateTime := time.NewTimer(AGGREGATOR_ITER_TIMEOUT)
	for {
		select {
		case <-psqlUpdateTime.C:
			tsNow := time.Now()
			request := fmt.Sprintf(
				consts.REQUEST_GET_STARTTIME,
				int(tsNow.Month()),
				int(tsNow.Weekday()),
				int(tsNow.Hour()),
				int(tsNow.Minute()))
			rows, err := o.postgres.DB.Queryx(request)
			o.log.Infof("Rows comming")
			if err != nil {
				o.log.Error(err)
			}

			o.operatorWG.Add(1)
			go func() {
				defer o.operatorWG.Done()
				o.processTask(rows)
			}()

			timeout := time.Duration(60 - (time.Now().Unix() - tsNow.Unix()))
			psqlUpdateTime = time.NewTimer(timeout * time.Second)
		}
	}
}

func (o *Operator) processTask(rows *sqlx.Rows) {
	var task structs.Task
	for rows.Next() {
		err := rows.StructScan(&task)
		if err != nil {
			o.log.Error(err)
			return
		}
		dateFolder := time.Now().Format(DATATIME_FORMAT)
		task.DumpFolder = fmt.Sprintf("'%s/%s'", task.DumpFolder, dateFolder)
		o.tasksChan <- task
	}

	if rows.Err() != nil {
		o.log.Error(rows.Err())
	}
}

func (o *Operator) getMessagesByDatabaseType(task structs.Task) (messages []structs.Task, err error) {
	switch task.DBType {
	case MONGO:
		messages, err = mongodb.GetMessages(task)
		if err != nil {
			return
		}
	default:
		err := fmt.Errorf("Driver for Database %s not found", task.DBType)
		o.stat.SendStatMessage(structs.FailOperator, task.UserID, task.DBID, task.TaskID, err)
	}
	return
}
