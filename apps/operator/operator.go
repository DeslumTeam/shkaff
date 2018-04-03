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
	logging "github.com/op/go-logging"
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

func InitOperator() (oper *Operator) {
	var err error
	oper = &Operator{
		postgres:  maindb.InitPSQL(),
		rabbit:    producer.InitAMQPProducer("mongodb"),
		tasksChan: make(chan structs.Task),
		log:       logger.GetLogs("Operator"),
		stat:      statsender.Run(),
	}
	oper.taskCache, err = cache.InitCacheDB()
	if err != nil {
		oper.log.Fatal(err)
	}
	return
}

func (oper *Operator) Run() {
	oper.operatorWG = sync.WaitGroup{}
	oper.operatorWG.Add(2)
	go oper.aggregator()
	go oper.taskSender()
	oper.log.Info("Start Operator")
	oper.operatorWG.Wait()
}

func (oper *Operator) Stop() {
	for i := 0; i < 2; i++ {
		oper.operatorWG.Done()
	}
	oper.log.Info("Stop Operator")
}

func (oper *Operator) taskSender() {
	var messages []structs.Task
	rabbit := oper.rabbit
	for task := range oper.tasksChan {
		switch dbType := task.DBType; dbType {
		case "mongodb":
			messages = mongodb.GetMessages(task)
		default:
			err := fmt.Sprintf("Driver for Database %s not found", task.DBType)
			oper.log.Info(err)
			continue
		}
		for _, msg := range messages {
			body, err := json.Marshal(msg)
			if err != nil {
				oper.log.Error(err)
				continue
			}
			if err := rabbit.Publish(body); err != nil {
				oper.log.Error(err)
				continue
			}
		}
	}
}

func (oper *Operator) aggregator() {
	psqlUpdateTime := time.NewTimer(10 * time.Second)
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
			rows, err := oper.postgres.DB.Queryx(request)
			if err != nil {
				oper.log.Error(err)
			} else {
				go oper.processTask(rows)
			}
			timeout := time.Duration(60 - (time.Now().Unix() - tsNow.Unix()))
			psqlUpdateTime = time.NewTimer(timeout * time.Second)
		}
	}
}

func (oper *Operator) processTask(rows *sqlx.Rows) {
	var task = structs.Task{}
	for rows.Next() {
		err := rows.StructScan(&task)
		if err != nil {
			oper.log.Error(err)
			return
		}
		dateFolder := time.Now().Format("2006-01-02_15_03")
		task.DumpFolder = fmt.Sprintf("%s/%s", task.DumpFolder, dateFolder)
		oper.tasksChan <- task
	}
}
