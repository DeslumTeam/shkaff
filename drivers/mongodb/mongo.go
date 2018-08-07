package mongodb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/DeslumTeam/shkaff/internal/structs"
	"gopkg.in/mgo.v2"
)

const (
	LOGINDBNAME         = "admin"         // Admin database name
	EMPTY_TASK          = "{}"            // Check empty task pattern
	MONGO_DIAL_TIUMEOUT = 5 * time.Second // Dial timeout seconds
)

func getAllDatabases(task structs.Task) (messages []structs.Task, err error) {
	url := fmt.Sprintf("%s:%d", task.Host, task.Port)
	session, err := mgo.DialWithTimeout(url, MONGO_DIAL_TIUMEOUT)
	if err != nil {
		return
	}
	defer session.Close()

	if task.DBUser != "" && task.DBPassword != "" {
		admindb := session.DB(LOGINDBNAME)
		err = admindb.Login(task.DBUser, task.DBPassword)
		if err != nil {
			return
		}
	}

	dbNames, err := session.DatabaseNames()
	if err != nil {
		return
	}

	for _, dbName := range dbNames {
		task.Database = dbName
		messages = append(messages, task)
	}

	return
}

func getCustomDatabases(task structs.Task) (messages []structs.Task, err error) {
	databases := make(map[string][]string)
	err = json.Unmarshal([]byte(task.Databases), &databases)
	if err != nil {
		return
	}

	for base := range databases {
		task.Database = base
		messages = append(messages, task)
	}

	return
}

// GetMessages get message for send on dump
func GetMessages(task structs.Task) (messages []structs.Task, err error) {
	if task.Databases == EMPTY_TASK {
		messages, err = getAllDatabases(task)
		return
	}

	messages, err = getCustomDatabases(task)
	return
}
