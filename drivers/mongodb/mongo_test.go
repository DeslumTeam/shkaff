package mongodb

import (
	"fmt"
	"shkaff/internal/structs"
	"testing"
)

func Test_mongoCliStruct_forEmptyDatabases(t *testing.T) {
	var task structs.Task
	task.Host = "127.0.0.1"
	task.Port = 27017
	var listTasks []structs.Task
	m := &mongoCliStruct{
		task:     task,
		messages: listTasks,
	}
	m.forEmptyDatabases()
	for _, x := range m.messages {
		fmt.Println(x.Database)
	}
	t.Error("1")
}
