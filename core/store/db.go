package store

import (
	"encoding/binary"
	"github.com/boltdb/bolt"
	"sync"
)

//Store interface
type Store interface {
	Open() error
	Read()
	Update(data []byte) error
	Close() error
	GetData() chan []byte
}

//DBStore...
type DBStore struct {
	Store
	db   *bolt.DB
	Data chan []byte
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

func (s *DBStore) Read() {
	s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("switch"))
		c := bucket.Cursor()
		for key, value := c.First(); key != nil; c.Next() {
			s.Data <- value
		}
		tx.DeleteBucket([]byte("switch"))
		return err
	})
}

func GetStore() Store {
	locker.Lock()
	defer locker.Unlock()
	if ds == nil {
		ds = new(DBStore)
		ds.Data = make(chan []byte)
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

func (s *DBStore) GetData() chan []byte {
	locker.Lock()
	defer locker.Unlock()
	return s.Data
}
