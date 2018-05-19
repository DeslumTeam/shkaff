package stat

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DeslumTeam/shkaff/internal/options"
	"github.com/DeslumTeam/shkaff/internal/structs"

	"github.com/jmoiron/sqlx"
	_ "github.com/kshvakov/clickhouse"
)

const (
	URI_TEMPLATE   = "tcp://%s:%d?debug=False"
	CHECKOUT_TIME  = 15
	INSERT_REQUEST = `
	INSERT INTO shkaff_stat (UserId, DbID, TaskId, Service, NewOperator,SuccessOperator,
		FailOperator, NewDump, SuccessDump, FailDump,
		NewRestore, SuccessRestore, FailRestore, Error, CreateDate)
	VALUES (:UserId, :DbID, :TaskId,:Service, :NewOperator, :SuccessOperator,
		:FailOperator, :NewDump, :SuccessDump, :FailDump,
		:NewRestore, :SuccessRestore, :FailRestore, :Error, :CreateDate)`

	SELECT_REQUEST = `
	SELECT
	 	sum(SuccessOperator) as Operator_Success,
		sum(FailOperator) as Operator_Fail,
		sum(SuccessDump) as Dump_Success,
		sum(FailDump) as Dump_Fail,
		sum(SuccessRestore) as Restore_Success,
		sum(FailRestore) as Restore_Fail
	FROM shkaff_stat
	WHERE (CreateDate == toDate(now()))`

	SELECT_ERRORS_REQUEST = `
	SELECT 
		Error as error,
		Service as service, 
		Count() as count 
	FROM shkaff_stat 
	WHERE (CreateDate == toDate(now())) AND (FailRestore<>0 or FailOperator<>0 or FailDump<>0)
	GROUP BY Error, Service`
)

type Errors struct {
	Error   string `db:"error"`
	Service string `db:"service"`
	Count   int    `db:"count"`
}

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
		_, err = stmt.Exec(sm.UserID, sm.DbID, sm.TaskID, sm.Service, sm.NewOperator, sm.SuccessOperator,
			sm.FailOperator, sm.NewDump, sm.SuccessDump, sm.FailDump, sm.NewRestore,
			sm.SuccessRestore, sm.FailRestore, sm.Error, sm.CreateDate)
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
func (s *StatDB) StandartStatSelect() (result map[string]map[string]interface{}, err error) {
	var row *sqlx.Row
	var columns []string
	row = s.DB.QueryRowx(SELECT_REQUEST)
	columns, err = row.Columns()
	if err != nil {
		return
	}
	res := make([]interface{}, len(columns))
	resP := make([]interface{}, len(columns))
	for i := range columns {
		resP[i] = &res[i]
	}
	err = row.Scan(resP...)
	if err != nil {
		return
	}
	result = make(map[string]map[string]interface{})
	for i, colName := range columns {
		val := resP[i].(*interface{})
		names := strings.Split(colName, "_")
		if len(names) != 2 {
			continue
		}
		service := names[0]
		status := names[1]
		if result[service] == nil {
			result[service] = make(map[string]interface{})
		}
		result[service][status] = *val
	}
	return
}

func (s *StatDB) dropList() {
	s.statMessageList = []structs.StatMessage{}
}

func (s *StatDB) SelectDailyErrors() (result [][]string, err error) {
	var task Errors

	rows, err := s.DB.Queryx(SELECT_ERRORS_REQUEST)
	if err != nil {
		return
	}
	for rows.Next() {
		err := rows.StructScan(&task)
		if err != nil {
			continue
		}
		var arr []string
		arr = append(arr, task.Error)
		arr = append(arr, task.Service)
		arr = append(arr, strconv.Itoa(task.Count))
		result = append(result, arr)
	}
	return
}
