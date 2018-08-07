package producer

import (
	"testing"
)

func TestInitAMQPProducer(t *testing.T) {
	tests := []struct {
		name           string
		wantURI        string
		queueName      string
		publishMessage string
		wantInitErr    bool
		wantPublishErr bool
	}{
		// {
		// 	name:           "InitTestOk",
		// 	wantURI:        "amqp://shkaff:shkaff@localhost:5672/shkaff_workers",
		// 	queueName:      "mongodb",
		// 	publishMessage: "{\"test\":1}",
		// 	wantInitErr:    false,
		// 	wantPublishErr: false,
		// },
		// {
		// 	name:           "InitTestErr",
		// 	wantURI:        "amqp://shkaff:shkaff@localhost:5672/shkaff_workers",
		// 	queueName:      "",
		// 	publishMessage: "{\"test\":2}",
		// 	wantInitErr:    true,
		// 	wantPublishErr: false,
		// },
		{
			name:           "PublishErr",
			wantURI:        "amqp://shkaff:shkaff@localhost:5672/shkaff_workers",
			queueName:      "test",
			publishMessage: "{\"test\":}",
			wantInitErr:    false,
			wantPublishErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQp, err := InitAMQPProducer(tt.queueName)
			if (err != nil) != tt.wantInitErr {
				t.Errorf("TestInitAMQPConsumer return %v", err)
				return
			}

			if gotQp.uri != tt.wantURI {
				t.Errorf("Fail URI: %v", tt.wantURI)
				return
			}

			if gotQp.log == nil {
				t.Errorf("Fail logging gotQp.log is nil")
				return
			}

			err = gotQp.Publish([]byte(tt.publishMessage))
			if (err != nil) != tt.wantPublishErr {
				t.Errorf("TestInitAMQPConsumer return %v", err)
				return
			}
		})
	}
}
