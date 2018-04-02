package cache

import (
	"os"
	"reflect"
	"testing"

	"github.com/DeslumTeam/shkaff/internal/consts"

	"github.com/syndtr/goleveldb/leveldb"
)

func TestInitCacheDB(t *testing.T) {
	tests := []struct {
		name      string
		wantCache *Cache
		wantErr   bool
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := InitCacheDB()
			if err != nil {
				t.Error(err)
			}
			_, err = os.Stat(consts.CACHEPATH)
			if err != nil {
				t.Error(err)
			}
			if cache == nil {
				t.Error("Cache descriptor nil")
			}
		})
		os.RemoveAll(consts.CACHEPATH)
	}
}

func TestCache_SetKV(t *testing.T) {
	var ts int64
	type fields struct {
		DB *leveldb.DB
	}
	type args struct {
		userID int
		dbID   int
		taskID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				userID: 1,
				dbID:   1,
				taskID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := InitCacheDB()
			if err != nil {
				t.Error(err)
				return
			}
			err = cache.SetKV(tt.args.userID, tt.args.dbID, tt.args.taskID)
			if err != nil {
				t.Error(err)
				return
			}
			ts, err = cache.GetKV(tt.args.userID, tt.args.dbID, tt.args.taskID)
			if err != nil {
				t.Error(err)
				return
			}
			if reflect.TypeOf(ts).String() != "int64" {
				t.Error("Return timestamp not int64")
				return
			}
			res, err := cache.ExistKV(tt.args.userID, tt.args.dbID, tt.args.taskID)
			if err != nil {
				t.Error(err)
				return
			}
			if !res {
				t.Errorf("Key %d|%d|%d not exists", tt.args.userID, tt.args.dbID, tt.args.taskID)
			}
			err = cache.DeleteKV(tt.args.userID, tt.args.dbID, tt.args.taskID)
			if err != nil {
				t.Error(err)
				return
			}
			err = cache.DeleteKV(tt.args.userID, tt.args.dbID, tt.args.taskID)
			if err != nil {
				t.Error(err)
				return
			}
			res, err = cache.ExistKV(tt.args.userID, tt.args.dbID, tt.args.taskID)
			if err == nil {
				t.Error("Key Exists not deleted")
				return
			}
		})
		os.RemoveAll(consts.CACHEPATH)
	}
}
