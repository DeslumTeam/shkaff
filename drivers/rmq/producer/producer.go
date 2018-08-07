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

const (
	EMPTY_EXCHANGE   = ""
	RMQ_CONTENT_TYPE = "application/json"
)

type RMQ struct {
	uri        string
	queueName  string
	Channel    *amqp.Channel
	Connect    *amqp.Connection
	Publishing *amqp.Publishing
	log        *logging.Logger
}

func InitAMQPProducer(queueName string) (qp *RMQ, err error) {
	if queueName == "" {
		return nil, fmt.Errorf("QueueName is empty")
	}

	cfg := options.InitControlConfig()
	qp = &RMQ{
		uri: fmt.Sprintf(consts.RMQ_URI_TEMPLATE, cfg.RMQ_USER,
			cfg.RMQ_PASS,
			cfg.RMQ_HOST,
			cfg.RMQ_PORT,
			cfg.RMQ_VHOST),
		queueName: queueName,
		log:       logger.GetLogs("RMQ Producer"),
	}

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

		return
	}
}

func (qp *RMQ) Publish(body []byte) (err error) {
	qp.Publishing = &amqp.Publishing{
		ContentType: RMQ_CONTENT_TYPE,
		Body:        body,
	}

	err = qp.Channel.Publish(EMPTY_EXCHANGE, qp.queueName, false, false, *qp.Publishing)
	return
}
