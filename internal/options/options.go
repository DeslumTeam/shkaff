package options

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"shkaff/internal/consts"
	"shkaff/internal/logger"
	"sync"

	logging "github.com/op/go-logging"
)

var (
	cc          *ShkaffConfig
	mutexConfig sync.Mutex
)

type ShkaffConfig struct {
	SHKAFF_UI_HOST        string         `"json:SHKAFF_UI_HOST"`
	SHKAFF_UI_PORT        int            `"json:SHKAFF_UI_PORT"`
	RMQ_HOST              string         `json:"RMQ_HOST"`
	RMQ_PORT              int            `json:"RMQ_PORT"`
	RMQ_USER              string         `json:"RMQ_USER"`
	RMQ_PASS              string         `json:"RMQ_PASS"`
	RMQ_VHOST             string         `json:"RMQ_VHOST"`
	DATABASE_HOST         string         `json:"DATABASE_HOST"`
	DATABASE_PORT         int            `json:"DATABASE_PORT"`
	DATABASE_USER         string         `json:"DATABASE_USER"`
	DATABASE_PASS         string         `json:"DATABASE_PASS"`
	DATABASE_DB           string         `json:"DATABASE_DB"`
	DATABASE_SSL          bool           `json:"DATABASE_SSL"`
	STATBASE_HOST         string         `json:"STATBASE_HOST"`
	STATBASE_PORT         int            `json:"STATBASE_PORT"`
	STATBASE_USER         string         `json:"STATBASE_USER"`
	STATBASE_PASS         string         `json:"STATBASE_PASS"`
	MONGO_RESTORE_HOST    string         `json:"MONGO_RESTORE_HOST"`
	MONGO_RESTORE_PORT    int            `json:"MONGO_RESTORE_PORT"`
	REFRESH_DATABASE_SCAN int            `json:"REFRESH_DATABASE_SCAN"`
	WORKERS               map[string]int `json:"WORKERS_COUNT"`
	log                   *logging.Logger
}

func InitControlConfig() *ShkaffConfig {
	mutexConfig.Lock()
	defer mutexConfig.Unlock()
	if cc != nil {
		return cc
	}
	cc = &ShkaffConfig{}
	var file []byte
	var err error
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if file, err = ioutil.ReadFile(dir + "/" + consts.CONFIG_FILE); err != nil {
		log.Fatalln(err)
		return nil
	}
	if err := json.Unmarshal(file, &cc); err != nil {
		log.Fatalln(err)
		return nil
	}
	cc.log = logger.GetLogs("Config")
	cc.validate()
	return cc
}

func (cc *ShkaffConfig) validate() {
	if cc.SHKAFF_UI_HOST == "" {
		log.Fatalf("Invalid Shkaff UI Host %s", cc.SHKAFF_UI_HOST)
	}
	if cc.SHKAFF_UI_PORT < 1025 || cc.SHKAFF_UI_PORT > 65535 {
		log.Fatalf("Invalid Shkaff UI Port %d", cc.SHKAFF_UI_PORT)
	}

	if cc.DATABASE_HOST == "" {
		log.Printf(consts.INVALID_DATABASE_HOST, consts.DEFAULT_HOST)
		cc.DATABASE_HOST = consts.DEFAULT_HOST
	}
	if cc.DATABASE_PORT < 1025 || cc.DATABASE_PORT > 65535 {
		log.Printf(consts.INVALID_DATABASE_PORT, cc.DATABASE_PORT, consts.DEFAULT_DATABASE_PORT)
		cc.DATABASE_PORT = consts.DEFAULT_DATABASE_PORT
	}
	if cc.DATABASE_DB == "" {
		log.Printf(consts.INVALID_DATABASE_DB, consts.DEFAULT_DATABASE_DB)
		cc.DATABASE_DB = consts.DEFAULT_DATABASE_DB
	}
	if cc.DATABASE_USER == "" {
		log.Fatalln(consts.INVALID_DATABASE_USER)
	}
	if cc.DATABASE_PASS == "" {
		log.Fatalln(consts.INVALID_DATABASE_PASSWORD)
	}

	if cc.RMQ_HOST == "" {
		log.Printf(consts.INVALID_AMQP_HOST, consts.DEFAULT_HOST)
		cc.RMQ_HOST = consts.DEFAULT_HOST
	}
	if cc.RMQ_PORT < 1025 || cc.RMQ_PORT > 65535 {
		log.Printf(consts.INVALID_AMQP_PORT, cc.RMQ_PORT, consts.DEFAULT_AMQP_PORT)
		cc.RMQ_PORT = consts.DEFAULT_AMQP_PORT
	}
	if cc.RMQ_USER == "" {
		log.Fatalln(consts.INVALID_AMQP_USER)
	}
	if cc.RMQ_PASS == "" {
		log.Fatalln(consts.INVALID_AMQP_PASSWORD)
	}

	if cc.STATBASE_HOST == "" {
		log.Printf(consts.INVALID_STATDB_HOST, consts.DEFAULT_HOST)
		cc.STATBASE_HOST = consts.DEFAULT_HOST
	}
	if cc.STATBASE_PORT < 1025 || cc.STATBASE_PORT > 65535 {
		log.Printf(consts.INVALID_STATDB_PORT, cc.STATBASE_PORT, consts.DEFAULT_STATDB_PORT)
		cc.STATBASE_PORT = consts.DEFAULT_STATDB_PORT
	}

	if cc.STATBASE_HOST == "" {
		log.Printf(consts.INVALID_STATDB_HOST, consts.DEFAULT_HOST)
		cc.STATBASE_HOST = consts.DEFAULT_HOST
	}
	if cc.STATBASE_PORT < 1025 || cc.STATBASE_PORT > 65535 {
		log.Printf(consts.INVALID_STATDB_PORT, cc.STATBASE_PORT, consts.DEFAULT_STATDB_PORT)
		cc.STATBASE_PORT = consts.DEFAULT_STATDB_PORT
	}
	if cc.MONGO_RESTORE_HOST == "" {
		log.Printf(consts.INVALID_MONGO_RESTORE_HOST, consts.DEFAULT_HOST)
		cc.MONGO_RESTORE_HOST = consts.DEFAULT_HOST
	}
	if cc.MONGO_RESTORE_PORT < 1025 || cc.MONGO_RESTORE_PORT > 65535 {
		log.Printf(consts.INVALID_MONGO_RESTORE_PORT, cc.MONGO_RESTORE_PORT, consts.DEFAULT_MONGO_RESTORE_PORT)
		cc.MONGO_RESTORE_PORT = consts.DEFAULT_STATDB_PORT
	}

	if cc.REFRESH_DATABASE_SCAN == 0 {
		cc.REFRESH_DATABASE_SCAN = consts.DEFAULT_REFRESH_DATABASE_SCAN
	}
	for database, workersCount := range cc.WORKERS {
		if workersCount > 0 {
			switch database {
			case "mongodb":
				cc.WORKERS[database] = workersCount
				cc.log.Infof("%s WorkersCount %d", database, workersCount)
			case "postgresql":
				cc.WORKERS[database] = workersCount
				cc.log.Infof("%s WorkersCount %d", database, workersCount)
			default:
				cc.log.Fatalf("Unknown Database %s", database)
			}
		} else {
			cc.log.Infof("%s WorkersCount %d", database, workersCount)
			delete(cc.WORKERS, database)
		}
	}
	return
}
