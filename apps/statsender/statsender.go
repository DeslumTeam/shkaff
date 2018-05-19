package statsender

import (
	"encoding/json"
	"time"

	"github.com/DeslumTeam/shkaff/drivers/rmq/consumer"
	"github.com/DeslumTeam/shkaff/drivers/rmq/producer"
	"github.com/DeslumTeam/shkaff/drivers/stat"
	"github.com/DeslumTeam/shkaff/internal/logger"
	"github.com/DeslumTeam/shkaff/internal/structs"

	logging "github.com/op/go-logging"
)

type StatSender struct {
	sChan    chan structs.StatMessage
	producer *producer.RMQ
	log      *logging.Logger
}

func Run() (statSender *StatSender) {
	statSender = &StatSender{
		sChan:    make(chan structs.StatMessage),
		producer: producer.InitAMQPProducer("shkaff_stat"),
		log:      logger.GetLogs("StatSender"),
	}
	statSender.log.Info("Start StatSender")
	go statSender.statSender()
	return
}

func (statSender *StatSender) SendStatMessage(action structs.Action, userID, dbid, taskID int, err error) {
	var statMessage structs.StatMessage
	statMessage.UserID = uint16(userID)
	statMessage.DbID = uint16(dbid)
	statMessage.TaskID = uint16(taskID)
	statMessage.CreateDate = time.Now()
	switch action {
	case 0:
		statMessage.NewOperator = 1
	case 1:
		statMessage.SuccessOperator = 1
	case 2:
		statMessage.Service = "Operator"
		statMessage.FailOperator = 1
		statMessage.Error = err.Error()
	case 3:
		statMessage.NewDump = 1
	case 4:
		statMessage.SuccessDump = 1
	case 5:
		statMessage.Service = "Dump"
		statMessage.FailDump = 1
		statMessage.Error = err.Error()
	case 6:
		statMessage.NewRestore = 1
	case 7:
		statMessage.SuccessRestore = 1
	case 8:
		statMessage.Service = "Restore"
		statMessage.FailRestore = 1
		statMessage.Error = err.Error()
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
