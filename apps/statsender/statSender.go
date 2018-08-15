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
	log := logger.GetLogs("StatSender")
	producer, err := producer.InitAMQPProducer("shkaff_stat")
	if err != nil {
		log.Fatalf("StatSender error %v", err)
	}

	statSender = &StatSender{
		sChan:    make(chan structs.StatMessage),
		producer: producer,
		log:      log,
	}
	statSender.log.Info("Start StatSender")
	go func() {
		statSender.statSender()
	}()

	return
}

func (statSender *StatSender) SendStatMessage(action structs.Action, userID, dbid, taskID int, err error) {
	statMessage := structs.StatMessage{
		UserID:     uint16(userID),
		DbID:       uint16(dbid),
		TaskID:     uint16(taskID),
		CreateDate: time.Now(),
	}

	switch action {
	case structs.NewOperator:
		statMessage.NewOperator = 1
	case structs.SuccessOperator:
		statMessage.SuccessOperator = 1
	case structs.FailOperator:
		statMessage.Service = "Operator"
		statMessage.FailOperator = 1
		statMessage.Error = err.Error()
	case structs.NewDump:
		statMessage.NewDump = 1
	case structs.SuccessDump:
		statMessage.SuccessDump = 1
	case structs.FailDump:
		statMessage.Service = "Dump"
		statMessage.FailDump = 1
		statMessage.Error = err.Error()
	case structs.NewRestore:
		statMessage.NewRestore = 1
	case structs.SuccessRestore:
		statMessage.SuccessRestore = 1
	case structs.FailRestore:
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
	var statMessage *structs.StatMessage
	statSender.consumer.InitConnection("shkaff_stat")
	for message := range statSender.consumer.Msgs {
		err := json.Unmarshal(message.Body, &statMessage)
		if err != nil {
			statSender.log.Error(err)
			continue
		}

		if statMessage == nil {
			statSender.log.Errorf("StatMessage is %v. Body is %v", statMessage, message.Body)
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
