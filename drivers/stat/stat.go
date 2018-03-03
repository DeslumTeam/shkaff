package stat

import (
	"fmt"
	"log"
	"shkaff/internal/options"
	"shkaff/internal/structs"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/kshvakov/clickhouse"
)

const (
	URI_TEMPLATE   = "tcp://%s:%d?debug=False"
	CHECKOUT_TIME  = 15
	INSERT_REQUEST = `
	INSERT INTO shkaff_stat (UserId, DbID, TaskId, NewOperator,SuccessOperator,
		FailOperator, ErrorOperator, NewDump, SuccessDump, FailDump, ErrorDump,
		NewRestore, SuccessRestore, FailRestore, ErrorRestore, CreateDate)
	VALUES (:UserId, :DbID, :TaskId, :NewOperator, :SuccessOperator,
		:FailOperator, :ErrorOperator, :NewDump, :SuccessDump, :FailDump, :ErrorDump,
		:NewRestore, :SuccessRestore, :FailRestore, :ErrorRestore, :CreateDate)`

	SELECT_REQUEST = `
	SELECT UserId, DbID, TaskId, sum(NewOperator) as NewOperator, 
	 sum(SuccessOperator) as SuccessOperator, sum(FailOperator) as FailOperator,
	 sum(NewDump) as NewDump, sum(SuccessDump) as SuccessDump, sum(FailDump) as FailDump,
	 sum(NewRestore) as NewRestore, sum(SuccessRestore) as SuccessRestore,
	 sum(FailRestore) as FailRestore  
	FROM shkaff_stat GROUP BY UserId, DbID, TaskId`
)

type StatDB struct {
	mutex           sync.Mutex
	uri             string
	statMessageList []structs.StatMessage
	DB              *sqlx.DB
}

func InitStat() (s *StatDB) {
	cfg := options.InitControlConfig()
	var err error
	s = new(StatDB)
	s.mutex = sync.Mutex{}
	s.uri = fmt.Sprintf(URI_TEMPLATE, cfg.STATBASE_HOST, cfg.STATBASE_PORT)
	for {
		s.DB, err = sqlx.Open("clickhouse", s.uri)
		if err == nil {
			break
		}
		log.Printf("ClickHouse: %s not connected\n", s.uri)
		time.Sleep(time.Second * 5)
	}

	go s.checkout()
	return
}

func (s *StatDB) Insert(statMessage structs.StatMessage) (err error) {
	s.mutex.Lock()
	s.statMessageList = append(s.statMessageList, statMessage)
	s.mutex.Unlock()
	return
}

func (s *StatDB) checkout() {
	for {
		timer := time.NewTimer(time.Second * CHECKOUT_TIME)
		select {
		case <-timer.C:
			if len(s.statMessageList) > 0 {
				s.inserBulk()
			}
		}
	}
}

func (s *StatDB) inserBulk() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	tx, err := s.DB.Begin()
	if err != nil {
		log.Println(err)
	}
	stmt, err := tx.Prepare(INSERT_REQUEST)
	if err != nil {
		log.Println(err)
	}
	for _, sm := range s.statMessageList {
		_, err = stmt.Exec(sm.UserID, sm.DbID, sm.TaskID, sm.NewOperator, sm.SuccessOperator,
			sm.FailOperator, sm.ErrorOperator, sm.NewDump, sm.SuccessDump, sm.FailDump,
			sm.ErrorDump, sm.NewRestore, sm.SuccessRestore, sm.FailRestore,
			sm.ErrorRestore, sm.CreateDate)
		if err != nil {
			log.Println(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}
	s.dropList()
}

//TODO Refactoring Very Ugly
func (s *StatDB) StandartStatSelect() (result map[string]interface{}, err error) {
	var row *sqlx.Row
	var columns []string
	row = s.DB.QueryRowx(SELECT_REQUEST)
	columns, err = row.Columns()
	if err != nil {
		return
	}
	res := make([]interface{}, len(columns))
	resP := make([]interface{}, len(columns))
	for i, _ := range columns {
		resP[i] = &res[i]
	}
	err = row.Scan(resP...)
	if err != nil {
		return
	}
	result = make(map[string]interface{})
	for i, colName := range columns {
		val := resP[i].(*interface{})
		result[colName] = *val
	}
	return
}

func (s *StatDB) dropList() {
	s.statMessageList = []structs.StatMessage{}
}
