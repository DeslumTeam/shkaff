package consumer

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
	Channel    *amqp.Channel
	Connect    *amqp.Connection
	Publishing *amqp.Publishing
	Msgs       <-chan amqp.Delivery
	log        *logging.Logger
}

func InitAMQPConsumer() (qp *RMQ) {
	cfg := options.InitControlConfig()
	qp = new(RMQ)
	qp.uri = fmt.Sprintf(consts.RMQ_URI_TEMPLATE, cfg.RMQ_USER,
		cfg.RMQ_PASS,
		cfg.RMQ_HOST,
		cfg.RMQ_PORT,
		cfg.RMQ_VHOST)
	qp.log = logger.GetLogs("RMQ Consumer")
	return
}

func (qp *RMQ) InitConnection(queueName string) {
	var err error
	if queueName == "" {
		qp.log.Fatal("Consumer queue name empty")
	}
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
	q, err := qp.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		qp.log.Fatal(err, "Failed to declare a queue")
	}
	if err = qp.Channel.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		qp.log.Fatal(err, "Failed to set QoS")
	}
	if msgs, err := qp.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	); err != nil {
		qp.log.Fatal(err, "Failed to register a consumer")
	} else {
		qp.Msgs = msgs
	}

}
