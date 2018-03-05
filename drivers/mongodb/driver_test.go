package mongodb

import (
	"testing"

	"github.com/DeslumTeam/shkaff/internal/databases"
	"github.com/DeslumTeam/shkaff/internal/structs"
)

func TestInitDriver(t *testing.T) {
	tests := []struct {
		name   string
		wantMp databases.DatabaseDriver
	}{
		{
			name: "InitMongoDB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drv := InitDriver()
			if drv == nil {
				t.Error("InitDriver() return nil")
			}
		})
	}
}

func TestMongoParams_ParamsToDumpString(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &MongoParams{
				cfg:      tt.fields.cfg,
				host:     tt.fields.host,
				port:     tt.fields.port,
				login:    tt.fields.login,
				password: tt.fields.password,
				ipv6:     tt.fields.ipv6,
				database: tt.fields.database,
				gzip:     tt.fields.gzip,
				parallelCollectionsNum: tt.fields.parallelCollectionsNum,
				dumpFolder:             tt.fields.dumpFolder,
				resultChan:             tt.fields.resultChan,
			}
			result := mp.ParamsToDumpString()
			if result != tt.wantDumpString {
				t.Errorf("Name: %s. Result %s. Want %s", tt.name, result, tt.wantDumpString)
				return
			}
		})
	}
}

func TestMongoParams_ParamsToRestoreString(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &MongoParams{
				cfg:      tt.fields.cfg,
				host:     tt.fields.host,
				port:     tt.fields.port,
				login:    tt.fields.login,
				password: tt.fields.password,
				ipv6:     tt.fields.ipv6,
				database: tt.fields.database,
				gzip:     tt.fields.gzip,
				parallelCollectionsNum: tt.fields.parallelCollectionsNum,
				dumpFolder:             tt.fields.dumpFolder,
				resultChan:             tt.fields.resultChan,
			}
			result := mp.ParamsToRestoreString()
			if result != tt.wantRestoreString {
				t.Errorf("Name: %s. Result %s. Want %s", tt.name, result, tt.wantRestoreString)
				return
			}
		})
	}
}

func TestMongoParams_setDBSettings(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &MongoParams{
				cfg:      tt.fields.cfg,
				host:     tt.fields.host,
				port:     tt.fields.port,
				login:    tt.fields.login,
				password: tt.fields.password,
				ipv6:     tt.fields.ipv6,
				database: tt.fields.database,
				gzip:     tt.fields.gzip,
				parallelCollectionsNum: tt.fields.parallelCollectionsNum,
				dumpFolder:             tt.fields.dumpFolder,
				resultChan:             tt.fields.resultChan,
			}
			task := new(structs.Task)
			task = &structs.Task{
				DBSettingsType: 1,
				DBSettingsID:   1,
				IsActive:       true,
				TaskID:         1,
				TaskName:       "CustomTest",
				DBID:           1,
				UserID:         1,
				Databases:      "",
				DBType:         "mongodb",
				Verb:           4,
				ThreadCount:    10,
				Gzip:           true,
				Ipv6:           true,
				Host:           "127.0.0.1",
				Port:           27017,
				DBUser:         "test",
				DBPassword:     "test",
				DumpFolder:     "/opt/dump",
				Database:       "test",
				IsDelete:       false,
			}
			mp.setDBSettings(task)
			if mp.host != task.Host {
				t.Errorf("Result %s. Want %s", mp.host, task.Host)
			}
			if mp.port != task.Port {
				t.Errorf("Result %d. Want %d", mp.port, task.Port)
			}
			if mp.ipv6 != task.Ipv6 {
				t.Errorf("Result %t. Want %t", mp.ipv6, task.Ipv6)
			}
			if mp.gzip != task.Gzip {
				t.Errorf("Result %t. Want %t", mp.gzip, task.Gzip)
			}
			if mp.login != task.DBUser {
				t.Errorf("Result %s. Want %s", mp.login, task.DBUser)
			}
			if mp.password != task.DBPassword {
				t.Errorf("Result %s. Want %s", mp.password, task.DBPassword)
			}
		})
	}
}

func TestMongoParams_Dump(t *testing.T) {
	for _, tt := range tests_tasks {
		t.Run(tt.name, func(t *testing.T) {
			mp := InitDriver()
			if err := mp.Dump(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("MongoParams.Dump() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
