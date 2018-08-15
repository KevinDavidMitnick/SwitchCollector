package store

import (
	"encoding/binary"
	"encoding/json"
	"github.com/boltdb/bolt"
	"os"
	"sync"
)

//Store interface
type Store interface {
	Open() error
	Read() string
	Update(data []byte) error
	Close() error
	CleanStale(timestamp int64)
}

//DBStore...
type DBStore struct {
	Store
	db *bolt.DB
	sync.RWMutex
}

var (
	locker      sync.RWMutex
	ds          *DBStore
	storeSwitch bool = true
)

//Open...
func (s *DBStore) Open() error {
	s.Lock()
	defer s.Unlock()
	var err error
	if s.db == nil {
		s.db, err = bolt.Open("opsultra.db", 0600, nil)
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("switch"))
		return err
	})
	return err
}

func (s *DBStore) Close() error {
	var err error
	if s.db != nil {
		err = s.db.Close()
		s.db = nil
	}
	return err
}

func (s *DBStore) Update(data []byte) error {
	s.Lock()
	defer s.Unlock()
	var err error
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("switch"))
		id, _ := bucket.NextSequence()
		return bucket.Put(itob(int(id)), data)
	})
	return err
}
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (s *DBStore) Read() string {
	s.Lock()
	defer s.Unlock()
	var data string
	s.db.Update(func(tx *bolt.Tx) error {
		bucket, _ := tx.CreateBucketIfNotExists([]byte("switch"))
		c := bucket.Cursor()
		key, value := c.First()
		if key != nil {
			data = string(value)
			return c.Delete()
		}
		return nil
	})
	return data
}

func (s *DBStore) CleanStale(timestamp int64) {
	s.Lock()
	defer s.Unlock()
	var flag bool = true
	data := make([]map[string]interface{}, 0)
	s.db.Update(func(tx *bolt.Tx) error {
		bucket, _ := tx.CreateBucketIfNotExists([]byte("switch"))
		c := bucket.Cursor()
		for key, value := c.First(); key != nil; key, value = c.Next() {
			if err := json.Unmarshal(value, &data); err == nil {
				if len(data) > 0 && data[0]["timestamp"].(float64) < float64(timestamp) {
					c.Delete()
				}
			}
		}
		if key, _ := bucket.Cursor().First(); key != nil {
			flag = false
		}
		return nil
	})
	if flag {
		s.Close()
		os.Remove("opsultra.db")
	}
}

func GetStore() Store {
	locker.Lock()
	defer locker.Unlock()
	if ds == nil {
		ds = &DBStore{}
		ds.Open()
	}
	if ds.db == nil {
		ds.Open()
	}
	return ds
}

func GetStoreStatus() bool {
	locker.Lock()
	defer locker.Unlock()
	return storeSwitch
}

func UpdateStoreStatus(status bool) {
	locker.Lock()
	defer locker.Unlock()
	storeSwitch = status
}
