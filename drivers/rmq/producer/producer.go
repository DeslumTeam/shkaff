package producer

import (
	"fmt"
	"time"

	"github.com/DeslumTeam/shkaff/internal/consts"
	"github.com/DeslumTeam/shkaff/internal/logger"
	"github.com/DeslumTeam/shkaff/internal/options"

	logging "github.com/op/go-logging"
	"github.com/streadway/amqp"
)

type RMQ struct {
	uri        string
	queueName  string
	Channel    *amqp.Channel
	Connect    *amqp.Connection
	Publishing *amqp.Publishing
	log        *logging.Logger
}

func InitAMQPProducer(queueName string) (qp *RMQ) {
	var err error
	cfg := options.InitControlConfig()
	qp = new(RMQ)
	qp.uri = fmt.Sprintf(consts.RMQ_URI_TEMPLATE, cfg.RMQ_USER,
		cfg.RMQ_PASS,
		cfg.RMQ_HOST,
		cfg.RMQ_PORT,
		cfg.RMQ_VHOST)
	qp.queueName = queueName
	qp.log = logger.GetLogs("RMQ Producer")
	for {
		qp.Connect, err = amqp.Dial(qp.uri)
		if err != nil {
			qp.log.Errorf("RMQ: %s not connected\n", qp.uri)
			time.Sleep(time.Second * 5)
			continue
		}
		qp.Channel, err = qp.Connect.Channel()
		if err != nil {
			qp.log.Errorf("Channel error %s", err)
			time.Sleep(time.Second * 5)
			continue
		}
		_, err = qp.Channel.QueueDeclare(
			qp.queueName, // name
			true,         // durable
			false,        // delete when unused
			false,        // exclusive
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			qp.log.Errorf("Queue declare error %s", err)
			time.Sleep(time.Second * 5)
			continue
		}
		qp.Publishing = new(amqp.Publishing)
		qp.Publishing.ContentType = "application/json"
		return
	}
}

func (qp *RMQ) Publish(body []byte) (err error) {
	qp.Publishing.Body = body
	if err = qp.Channel.Publish("", qp.queueName, false, false, *qp.Publishing); err != nil {
		return
	}
	return
}
