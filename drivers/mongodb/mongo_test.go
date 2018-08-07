package mongodb

import (
	"testing"

	"github.com/DeslumTeam/shkaff/internal/structs"
)

func Test_getAllDatabases(t *testing.T) {
	type args struct {
		task structs.Task
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Default empty database",
			args: args{
				task: structs.Task{
					DBSettingsType: 1,
					DBSettingsID:   1,
					IsActive:       true,
					TaskID:         1,
					TaskName:       "test",
					DBID:           1,
					UserID:         1,
					Databases:      "{}",
					DBType:         "mongodb",
					Verb:           1,
					ThreadCount:    4,
					Gzip:           false,
					Ipv6:           false,
					Host:           "localhost",
					Port:           27017,
					DBUser:         "shkaff",
					DBPassword:     "shkaff",
					DumpFolder:     "for_test",
					Database:       "admin",
					IsDelete:       false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := getAllDatabases(tt.args.task)
			if err != nil {
				t.Errorf("getAllDatabases() error = %v", err)
				return
			}

			if len(tasks) != 2 {
				t.Errorf("getAllDatabases() task count !=2  result : %v", len(tasks))
				return
			}

			for _, task := range tasks {
				if task.Database == "" {
					t.Errorf("getAllDatabases() task.Database is empty")
				}

				if task.Database != "local" && task.Database != "admin" {
					t.Errorf("getAllDatabases() task.Database is not valid")
				}
			}
		})
	}
}

func Test_getCustomDatabases(t *testing.T) {
	type args struct {
		task structs.Task
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Default empty database",
			args: args{
				task: structs.Task{
					DBSettingsType: 1,
					DBSettingsID:   1,
					IsActive:       true,
					TaskID:         1,
					TaskName:       "test",
					DBID:           1,
					UserID:         1,
					Databases:      "{\"test1\": [\"db_1\"], \"test2\": [\"db_2\"], \"test3\": [\"db_3\"]}",
					DBType:         "mongodb",
					Verb:           1,
					ThreadCount:    4,
					Gzip:           false,
					Ipv6:           false,
					Host:           "localhost",
					Port:           27017,
					DBUser:         "shkaff",
					DBPassword:     "shkaff",
					DumpFolder:     "for_test",
					Database:       "",
					IsDelete:       false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := getCustomDatabases(tt.args.task)
			if err != nil {
				t.Errorf("getCustomDatabases() error = %v", err)
				return
			}

			if len(tasks) != 3 {
				t.Errorf("getCustomDatabases() task count !=3  result : %v", len(tasks))
				return
			}

			for _, task := range tasks {
				if task.Database == "" {
					t.Errorf("getCustomDatabases() task.Database is empty")
				}

				if task.Database == "test1" {
					continue
				}

				if task.Database == "test2" {
					continue
				}

				if task.Database == "test3" {
					continue
				}

				t.Errorf("getCustomDatabases() task.Database is not valid. Result: %v", task.Database)
			}
		})
	}
}

func TestGetMessages(t *testing.T) {
	type args struct {
		task structs.Task
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "All",
			args: args{
				task: structs.Task{
					DBSettingsType: 1,
					DBSettingsID:   1,
					IsActive:       true,
					TaskID:         1,
					TaskName:       "test",
					DBID:           1,
					UserID:         1,
					Databases:      "{}",
					DBType:         "mongodb",
					Verb:           1,
					ThreadCount:    4,
					Gzip:           false,
					Ipv6:           false,
					Host:           "localhost",
					Port:           27017,
					DBUser:         "shkaff",
					DBPassword:     "shkaff",
					DumpFolder:     "for_test",
					Database:       "admin",
					IsDelete:       false,
				},
			},
		},
		{
			name: "Custom",
			args: args{
				task: structs.Task{
					DBSettingsType: 1,
					DBSettingsID:   1,
					IsActive:       true,
					TaskID:         1,
					TaskName:       "test",
					DBID:           1,
					UserID:         1,
					Databases:      "{\"test1\": [\"db_1\"], \"test2\": [\"db_2\"], \"test3\": [\"db_3\"]}",
					DBType:         "mongodb",
					Verb:           1,
					ThreadCount:    4,
					Gzip:           false,
					Ipv6:           false,
					Host:           "localhost",
					Port:           27017,
					DBUser:         "shkaff",
					DBPassword:     "shkaff",
					DumpFolder:     "for_test",
					Database:       "",
					IsDelete:       false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := GetMessages(tt.args.task)
			if err != nil {
				t.Errorf("GetMessages() error = %v", err)
				return
			}

			if len(tasks) == 0 {
				t.Errorf("GetMessages getAllDatabases() task count == 0")
				return
			}

			switch tt.name {
			case "All":
				if len(tasks) != 2 {
					t.Errorf("GetMessages getAllDatabases() task count !=2  result : %v", len(tasks))
					return
				}

				for _, task := range tasks {
					if task.Database == "" {
						t.Errorf("GetMessages getAllDatabases() task.Database is empty")
					}

					if task.Database != "local" && task.Database != "admin" {
						t.Errorf("GetMessages getAllDatabases() task.Database is not valid")
					}
				}
			case "Custom":
				for _, task := range tasks {
					if task.Database == "" {
						t.Errorf("GetMessages getCustomDatabases() task.Database is empty")
					}

					if task.Database == "test1" {
						continue
					}

					if task.Database == "test2" {
						continue
					}

					if task.Database == "test3" {
						continue
					}

					t.Errorf("GetMessages getCustomDatabases() task.Database is not valid. Result: %v", task.Database)
				}
			}
		})
	}
}
