package store

import (
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
	"sync"
)

//Store interface
type Store interface {
	Open() error
	Read() ([]byte, error)
	Update(data []byte) error
	Close() error
}

//DBStore...
type DBStore struct {
	Store
	db *bolt.DB
}

var (
	locker      sync.RWMutex
	ds          *DBStore
	storeSwitch bool
)

//Open...
func (s *DBStore) Open() error {
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

func (s *DBStore) Read() ([]byte, error) {
	var err error
	var data []byte
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("switch"))
		c := bucket.Cursor()
		key, value := c.First()
		if key == nil {
			return fmt.Errorf("read empty bucket")
		}
		data = value
		if err := bucket.Delete(key); err == nil {
			fmt.Println("delete,value:", string(value), "success!")
		} else {
			fmt.Println("delete,value:", string(value), "failure!")
		}
		return nil
	})
	return data, err
}

func GetStore() Store {
	locker.Lock()
	defer locker.Unlock()
	if ds == nil {
		ds = new(DBStore)
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
