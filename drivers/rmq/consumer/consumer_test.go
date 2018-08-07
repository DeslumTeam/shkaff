package consumer

import (
	"testing"
)

func TestInitAMQPConsumer(t *testing.T) {
	tests := []struct {
		name      string
		wantURI   string
		queueName string
		wantErr   bool
	}{
		{
			name:      "InitTest",
			wantURI:   "amqp://shkaff:shkaff@localhost:5672/shkaff_workers",
			queueName: "mongodb",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQp := InitAMQPConsumer()
			if gotQp.uri != tt.wantURI {
				t.Errorf("Fail URI: %v", tt.wantURI)
				return
			}

			if gotQp.log == nil {
				t.Errorf("Fail logging gotQp.log is nil")
				return
			}

			err := gotQp.InitConnection(tt.queueName)
			if err != nil && !tt.wantErr {
				t.Errorf("TestInitAMQPConsumer return %v", err)
				return
			}

		})
	}
}
