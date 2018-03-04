package mongodb

import (
	"github.com/DeslumTeam/shkaff/internal/structs"
)

type argsTaskCases struct {
	task *structs.Task
}

var tests_tasks = []struct {
	name    string
	args    argsTaskCases
	wantErr bool
}{
	{
		name: "0",
		args: argsTaskCases{
			task: &structs.Task{
				Host:     "0.0.0.0",
				Port:     27017,
				Database: "admin",
			},
		},
	},
	{
		name: "1",
		args: argsTaskCases{
			task: &structs.Task{
				Host:     "0.0.0.0",
				Port:     27017,
				Database: "admin",
				Gzip:     true,
			},
		},
	},
	{
		name: "2",
		args: argsTaskCases{
			task: &structs.Task{
				Host:        "0.0.0.0",
				Port:        27017,
				Database:    "admin",
				ThreadCount: 10,
			},
		},
	},
	{
		name: "3",
		args: argsTaskCases{
			task: &structs.Task{
				Host:        "0.0.0.0",
				Port:        27017,
				Database:    "admin",
				ThreadCount: 10,
			},
		},
	},
	{
		name: "4",
		args: argsTaskCases{
			task: &structs.Task{
				Database:    "admin",
				ThreadCount: 10,
				DumpFolder:  "../",
			},
		},
	},
	{
		name: "5",
		args: argsTaskCases{
			task: &structs.Task{},
		},
	},
}
