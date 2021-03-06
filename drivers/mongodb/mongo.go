package mongodb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/DeslumTeam/shkaff/internal/logger"
	"github.com/DeslumTeam/shkaff/internal/structs"

	logging "github.com/op/go-logging"
	"gopkg.in/mgo.v2"
)

type mongoCliStruct struct {
	task     structs.Task
	messages []structs.Task
	log      *logging.Logger
}

func (m *mongoCliStruct) emptyDB() (err error) {
	url := fmt.Sprintf("%s:%d", m.task.Host, m.task.Port)

	session, err := mgo.DialWithTimeout(url, 5*time.Second)
	if err != nil {
		m.log.Error(err)
		return
	}
	defer session.Close()

	if m.task.DBUser != "" && m.task.DBPassword != "" {
		admindb := session.DB("admin")
		err = admindb.Login(m.task.DBUser, m.task.DBPassword)
		if err != nil {
			m.log.Errorf("+++ %v \n", err)
			return
		}
	}
	dbNames, err := session.DatabaseNames()
	if err != nil {
		m.log.Error(err)
		return
	}
	for _, dbName := range dbNames {
		m.task.Database = dbName
		m.messages = append(m.messages, m.task)
	}
	return
}

func (m *mongoCliStruct) fillDB() (err error) {
	databases := make(map[string][]string)
	err = json.Unmarshal([]byte(m.task.Databases), &databases)
	if err != nil {
		m.log.Error("Error unmarshal databases", databases, err)
		return
	}
	for base := range databases {
		m.task.Database = base
		m.messages = append(m.messages, m.task)
	}
	return
}

func GetMessages(task structs.Task) (err error, caches []structs.Task) {
	var mongo = new(mongoCliStruct)
	mongo.task = task
	mongo.log = logger.GetLogs("Mongo")
	if task.Databases == "{}" {
		err = mongo.emptyDB()
	} else {
		err = mongo.fillDB()
	}
	return err, mongo.messages

}
