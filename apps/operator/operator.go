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

func InitOperator() (o *Operator) {
	var err error
	o = &Operator{
		postgres:  maindb.InitPSQL(),
		rabbit:    producer.InitAMQPProducer("mongodb"),
		tasksChan: make(chan structs.Task),
		log:       logger.GetLogs("Operator"),
		stat:      statsender.Run(),
	}
	o.taskCache, err = cache.InitCacheDB()
	if err != nil {
		o.log.Fatal(err)
	}
	return
}

func (o *Operator) Run() {
	o.operatorWG = sync.WaitGroup{}
	o.operatorWG.Add(2)
	go o.aggregator()
	go o.taskSender()
	o.log.Info("Start Operator")
	o.operatorWG.Wait()
}

func (o *Operator) Stop() {
	for i := 0; i < 2; i++ {
		o.operatorWG.Done()
	}
	o.log.Info("Stop Operator")
}

func (o *Operator) taskSender() {
	var messages []structs.Task
	for task := range o.tasksChan {

		switch dbType := task.DBType; dbType {
		case "mongodb":
			messages = mongodb.GetMessages(task)
		default:
			err := fmt.Sprintf("Driver for Database %s not found", task.DBType)
			o.log.Info(err)
			continue
		}
		for _, msg := range messages {
			body, err := json.Marshal(msg)
			if err != nil {
				o.log.Error(err)
				continue
			}

			err = o.rabbit.Publish(body)
			if err != nil {
				o.log.Error(err)
				continue
			}
		}
	}
}

func (o *Operator) aggregator() {
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
			rows, err := o.postgres.DB.Queryx(request)
			if err != nil {
				o.log.Error(err)
			} else {
				go o.processTask(rows)
			}
			timeout := time.Duration(60 - (time.Now().Unix() - tsNow.Unix()))
			psqlUpdateTime = time.NewTimer(timeout * time.Second)
		}
	}
}

func (o *Operator) processTask(rows *sqlx.Rows) {
	var task = structs.Task{}
	for rows.Next() {
		err := rows.StructScan(&task)
		if err != nil {
			o.log.Error(err)
			return
		}
		dateFolder := time.Now().Format("2006-01-02")
		task.DumpFolder = fmt.Sprintf("%s/%s", task.DumpFolder, dateFolder)
		o.tasksChan <- task
	}
}
