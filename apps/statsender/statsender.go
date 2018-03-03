package statsender

import (
	"encoding/json"
	"shkaff/drivers/rmq/consumer"
	"shkaff/drivers/rmq/producer"
	"shkaff/drivers/stat"
	"shkaff/internal/logger"
	"shkaff/internal/structs"
	"time"

	logging "github.com/op/go-logging"
)

type StatSender struct {
	sChan    chan structs.StatMessage
	producer *producer.RMQ
	log      *logging.Logger
}

func Run() (statSender *StatSender) {
	statSender = new(StatSender)
	statSender.sChan = make(chan structs.StatMessage)
	statSender.producer = producer.InitAMQPProducer("shkaff_stat")
	statSender.log = logger.GetLogs("StatSender")
	statSender.log.Info("Start StatSender")
	go statSender.statSender()
	return
}

func (statSender *StatSender) SendStatMessage(action structs.Action, userID, dbid, taskID int, err error) {
	var statMessage structs.StatMessage
	statMessage.UserID = uint16(userID)
	statMessage.DbID = uint16(dbid)
	statMessage.TaskID = uint16(taskID)
	statMessage.CreateDate = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	switch action {
	case 0:
		statMessage.NewOperator = 1
	case 1:
		statMessage.SuccessOperator = 1
	case 2:
		statMessage.FailOperator = 1
		statMessage.ErrorOperator = err.Error()
	case 3:
		statMessage.NewDump = 1
	case 4:
		statMessage.SuccessDump = 1
	case 5:
		statMessage.FailDump = 1
		statMessage.ErrorDump = err.Error()
	case 6:
		statMessage.NewRestore = 1
	case 7:
		statMessage.SuccessRestore = 1
	case 8:
		statMessage.FailRestore = 1
		statMessage.ErrorRestore = err.Error()
	}
	statSender.sChan <- statMessage
}

func (statSender *StatSender) statSender() {
	for {
		select {
		case statMsg := <-statSender.sChan:
			msg, err := json.Marshal(statMsg)
			if err != nil {
				statSender.log.Error(err)
				continue
			}
			err = statSender.producer.Publish(msg)
			if err != nil {
				statSender.log.Error(err)
				continue
			}
		}
	}
}

//StatWorker

type statWorker struct {
	consumer *consumer.RMQ
	statDB   *stat.StatDB
	log      *logging.Logger
}

func InitStatSender() (sw *statWorker) {
	sw = new(statWorker)
	sw.statDB = stat.InitStat()
	sw.consumer = consumer.InitAMQPConsumer()
	sw.log = logger.GetLogs("StatWorker")
	return
}

func (statSender *statWorker) Run() {
	statSender.log.Info("Start StatWorker")
	var statMessage structs.StatMessage
	statSender.consumer.InitConnection("shkaff_stat")
	for message := range statSender.consumer.Msgs {
		err := json.Unmarshal(message.Body, &statMessage)
		if err != nil {
			statSender.log.Error(err)
			continue
		}

		err = statSender.statDB.Insert(statMessage)
		if err != nil {
			statSender.log.Error(err)
			continue
		}
		message.Ack(false)
	}
}

func (statSender *statWorker) Stop() {
	err := statSender.consumer.Channel.Close()
	if err != nil {
		statSender.log.Error(err)
	}
	statSender.log.Info("Stop StatWorker")

}
