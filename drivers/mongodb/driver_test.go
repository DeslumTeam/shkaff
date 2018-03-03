package mongodb

import (
	"shkaff/internal/databases"
	"testing"
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
			if result != tt.wantCommandString {
				t.Errorf("Name: %s. Result %s. Want %s", tt.name, result, tt.wantCommandString)
				return
			}
		})
	}
}

// func TestMongoParams_setDBSettings(t *testing.T) {
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mp := &MongoParams{
// 				cfg:      tt.fields.cfg,
// 				host:     tt.fields.host,
// 				port:     tt.fields.port,
// 				login:    tt.fields.login,
// 				password: tt.fields.password,
// 				ipv6:     tt.fields.ipv6,
// 				database: tt.fields.database,
// 				gzip:     tt.fields.gzip,
// 				parallelCollectionsNum: tt.fields.parallelCollectionsNum,
// 				dumpFolder:             tt.fields.dumpFolder,
// 				resultChan:             tt.fields.resultChan,
// 			}

// 			mp.setDBSettings(tt.args.task)
// 		})
// 	}
// }
