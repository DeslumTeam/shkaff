package producer

import (
	"fmt"
	"shkaff/internal/consts"
	"shkaff/internal/logger"
	"shkaff/internal/options"
	"time"

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
	cfg := options.InitControlConfig()
	qp = new(RMQ)
	qp.uri = fmt.Sprintf(consts.RMQ_URI_TEMPLATE, cfg.RMQ_USER,
		cfg.RMQ_PASS,
		cfg.RMQ_HOST,
		cfg.RMQ_PORT,
		cfg.RMQ_VHOST)
	qp.queueName = queueName
	qp.log = logger.GetLogs("RMQ Producer")
	qp.initConnection()
	return
}

func (qp *RMQ) initConnection() {
	var err error
	for {
		qp.Connect, err = amqp.Dial(qp.uri)
		if err == nil {
			break
		}
		qp.log.Errorf("RMQ: %s not connected\n", qp.uri)
		time.Sleep(time.Second * 5)
	}
	if qp.Channel, err = qp.Connect.Channel(); err != nil {
		qp.log.Fatal(err)
	}
	if _, err = qp.Channel.QueueDeclare(
		qp.queueName, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		qp.log.Fatal(err)
	}
	qp.Publishing = new(amqp.Publishing)
	qp.Publishing.ContentType = "application/json"
}

func (qp *RMQ) Publish(body []byte) (err error) {
	qp.Publishing.Body = body
	if err = qp.Channel.Publish("", qp.queueName, false, false, *qp.Publishing); err != nil {
		return
	}
	return
}
