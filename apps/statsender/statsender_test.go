package statsender

import (
	"sync"
	"testing"

	"github.com/DeslumTeam/shkaff/drivers/rmq/producer"
	"github.com/DeslumTeam/shkaff/internal/structs"
	logging "github.com/op/go-logging"
)

func TestStatSender_SendStatMessage(t *testing.T) {
	testChan := make(chan structs.StatMessage)
	type fields struct {
		sChan    chan structs.StatMessage
		producer *producer.RMQ
		log      *logging.Logger
	}
	type args struct {
		action structs.Action
		userID int
		dbid   int
		taskID int
		err    error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "first",
			args: args{
				action: 0,
				userID: 1,
				dbid:   1,
				taskID: 1,
			},
		},
	}
	statSender := &StatSender{
		sChan: testChan,
	}
	var wg sync.WaitGroup

	for {
		f, _ := <-statSender.sChan
		t.Log(f.CreateDate)
		wg.Done()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				statSender.SendStatMessage(tt.args.action, tt.args.userID, tt.args.dbid, tt.args.taskID, tt.args.err)
			})
		}
		close(statSender.sChan)
	}()
	wg.Wait()
}
