package mongodb

import (
	"encoding/json"
	"fmt"
	"shkaff/internal/logger"
	"shkaff/internal/structs"
	"time"

	logging "github.com/op/go-logging"
	"gopkg.in/mgo.v2"
)

type mongoCliStruct struct {
	task     structs.Task
	messages []structs.Task
	log      *logging.Logger
}

func (m *mongoCliStruct) emptyDB() {
	url := fmt.Sprintf("%s:%d", m.task.Host, m.task.Port)
	session, err := mgo.DialWithTimeout(url, 5*time.Second)
	if err != nil {
		m.log.Error(err)
		return
	}
	defer session.Close()

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

func (m *mongoCliStruct) fillDB() {
	databases := make(map[string][]string)
	err := json.Unmarshal([]byte(m.task.Databases), &databases)
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

func GetMessages(task structs.Task) (caches []structs.Task) {
	var mongo = new(mongoCliStruct)
	mongo.task = task
	mongo.log = logger.GetLogs("Mongo")
	if task.Databases == "{}" {
		mongo.emptyDB()
	} else {
		mongo.fillDB()
	}
	return mongo.messages

}
