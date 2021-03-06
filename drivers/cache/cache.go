package cache

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/DeslumTeam/shkaff/internal/consts"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	TTL = 60 * 60 * 24 // 1 day = 24 hours = 84600 seconds
)

type Cache struct {
	*leveldb.DB
}

func InitCacheDB() (cache *Cache, err error) {
	err = os.MkdirAll(consts.CACHEPATH, 0755)
	if err != nil {
		return nil, err
	}
	cache = new(Cache)
	cache.DB, err = leveldb.OpenFile(consts.CACHEPATH, nil)
	if err != nil {
		return nil, err
	}
	return
}

func (cache *Cache) SetKV(userID, dbID, taskID int) (err error) {
	key := []byte(fmt.Sprintf("%d|%d|%d", userID, dbID, taskID))
	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, uint64(time.Now().Unix()+TTL))
	err = cache.DB.Put(key, value, nil)
	if err != nil {
		return err
	}
	return
}

func (cache *Cache) GetKV(userID, dbID, taskID int) (timestamp int64, err error) {
	key := []byte(fmt.Sprintf("%d|%d|%d", userID, dbID, taskID))
	valB, err := cache.DB.Get(key, nil)
	if err != nil {
		return 0, err
	}
	timestamp = int64(binary.LittleEndian.Uint64(valB))
	if timestamp < time.Now().Unix() {
		return 0, errors.New("Expire key")
	}
	return timestamp, nil
}

func (cache *Cache) DeleteKV(userID, dbID, taskID int) (err error) {
	key := []byte(fmt.Sprintf("%d|%d|%d", userID, dbID, taskID))
	err = cache.DB.Delete(key, nil)
	if err != nil {
		return err
	}
	return
}

func (cache *Cache) ExistKV(userID, dbID, taskID int) (result bool, err error) {
	key := []byte(fmt.Sprintf("%d|%d|%d", userID, dbID, taskID))
	res, err := cache.Get(key, nil)
	if err == nil {
		timestamp := int64(binary.LittleEndian.Uint64(res))
		if timestamp < time.Now().Unix() {
			return false, nil
		}
	}
	if err != nil {
		return false, err
	}
	if res == nil {
		return false, nil
	}
	return true, nil
}
