package worker

import (
	"shkaff/apps/statsender"
	"shkaff/drivers/maindb"
	"shkaff/drivers/rmq/consumer"
	"shkaff/internal/databases"
	"testing"

	logging "github.com/op/go-logging"
)

func Test_worker_getDatabaseType_success(t *testing.T) {

	type fields struct {
		dumpChan     chan string
		databaseName string
		postgres     *maindb.PSQL
		workRabbit   *consumer.RMQ
		stat         *statsender.StatSender
		log          *logging.Logger
	}
	tests := []struct {
		name         string
		fields       fields
		wantDbDriver databases.DatabaseDriver
		wantErr      bool
	}{
		{
			name: "Test MongoDB request",
			fields: fields{
				databaseName: "mongodb",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &worker{
				dumpChan:     tt.fields.dumpChan,
				databaseName: tt.fields.databaseName,
				postgres:     tt.fields.postgres,
				workRabbit:   tt.fields.workRabbit,
				stat:         tt.fields.stat,
				log:          tt.fields.log,
			}
			gotDbDriver, err := w.getDatabaseType()
			if (err != nil) != tt.wantErr {
				t.Errorf("worker.getDatabaseType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDbDriver == nil {
				t.Errorf("DB driver Handle is nil")
				return
			}
		})
	}
}

func Test_worker_getDatabaseType_fail(t *testing.T) {

	type fields struct {
		dumpChan     chan string
		databaseName string
		postgres     *maindb.PSQL
		workRabbit   *consumer.RMQ
		stat         *statsender.StatSender
		log          *logging.Logger
	}
	tests := []struct {
		name         string
		fields       fields
		wantDbDriver databases.DatabaseDriver
		wantErr      bool
	}{
		{
			name: "Test Fail request",
			fields: fields{
				databaseName: "rabbitmq",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &worker{
				dumpChan:     tt.fields.dumpChan,
				databaseName: tt.fields.databaseName,
				postgres:     tt.fields.postgres,
				workRabbit:   tt.fields.workRabbit,
				stat:         tt.fields.stat,
				log:          tt.fields.log,
			}
			gotDbDriver, err := w.getDatabaseType()
			if err != nil && err.Error() != "Driver rabbitmq not found" {
				t.Errorf("worker.getDatabaseType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDbDriver != nil && err != nil {
				t.Errorf("DB driver Handle is not nil ")
				return
			}

		})
	}
}

func TestInitWorkerSuccess(t *testing.T) {
	tests := []struct {
		name   string
		wantWs *workersStarter
	}{
		{
			name: "SucessInitWorker",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWs := InitWorker()
			if gotWs == nil {
				t.Errorf("InitWorker() = %v, want handle", gotWs)
			}
			if gotWs.workerCount == 0 {
				t.Errorf("workerCount() = %d, want number greater 0", gotWs.workerCount)

			}
			wCount := len(gotWs.workers)
			if wCount != gotWs.workerCount {
				t.Errorf("workers length = %d, want number greater %d", wCount, gotWs.workerCount)
			}
		})
	}
}
