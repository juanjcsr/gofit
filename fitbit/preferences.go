package fitbit

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type Preferences struct {
	db         *bolt.DB
	bucketName string
}

const dbName = "gofit.db"
const bucketName = "prefs"

func (p *Preferences) Open() (bool, error) {
	var err error
	p.db, err = bolt.Open(dbName, 0600, nil)
	if err != nil {
		return false, fmt.Errorf("open database: %s", err)
	}
	p.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return true, nil
}

// Update saves the K/V to the db
func (p *Preferences) Update(key string, value string) {
	p.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		err := bucket.Put([]byte(key), []byte(value))
		//fmt.Printf("error in update: %v", err)
		return err
	})
}

// Read retrieves the value of the key
func (p *Preferences) Read(key string) string {
	var value string
	p.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		values := bucket.Get([]byte(key))
		//fmt.Printf("READING %s with value %s", bucketName, values)
		if values != nil {
			value = string(values)
		}
		return nil
	})
	return value
}

// Close closes the db and handles its shutdown
func (p *Preferences) Close() {
	p.db.Close()
}
