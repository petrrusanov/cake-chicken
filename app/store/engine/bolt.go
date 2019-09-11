package engine

import (
	"fmt"
	"github.com/petrrusanov/cake-chicken/app/store/models"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"strconv"
)

// Store extends engine interface with specific fields to the store
type Store struct {
	Interface

	db *bolt.DB
}

const (
	chickensBucket = "chickens"
	cakesBucket    = "cakes"
)

// NewBolt creates new store with specified params
func NewBolt(path string) (*Store, error) {
	db, err := bolt.Open(path, 0600, nil)

	if err != nil {
		return nil, errors.Wrap(err, "failed to open bolt db file")
	}

	result := Store{db: db}
	err = result.prepare()

	return &result, errors.Wrap(err, "failed to prepare bolt db")
}

// prepare collections with all indexes
func (s *Store) prepare() error {
	return nil
}

// Close bolt store
func (s *Store) Close() error {
	return s.db.Close()
}

// AddChicken Add to a chicken counter for a username
func (s *Store) AddChicken(username string, prefix string) (models.UserCounter, error) {
	bucketName := fmt.Sprintf("%s.%s", prefix, chickensBucket)

	var userCounter models.UserCounter

	err := s.ensureBucketExists(bucketName)

	if err != nil {
		return userCounter, err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		counterString := string(b.Get([]byte(username)))

		var counter int
		var err error

		if counterString == "" {
			counter = 0
		} else {
			counter, err = strconv.Atoi(counterString)

			if  err != nil {
				return errors.Wrap(err, "chicken counter is not a number")
			}
		}

		counter = counter + 1
		counterString = strconv.Itoa(counter)

		err = b.Put([]byte(username), []byte(counterString))

		if err == nil {
			userCounter = models.UserCounter{Username: username, Count: counter}
		}

		return err
	})

	return userCounter, err
}

// FulfillChicken Fulfill chicken promise for a username
func (s *Store) FulfillChicken(username string, prefix string) (models.UserCounter, error) {
	bucketName := fmt.Sprintf("%s.%s", prefix, chickensBucket)

	var userCounter models.UserCounter

	err := s.ensureBucketExists(bucketName)

	if err != nil {
		return userCounter, err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		counterString := string(b.Get([]byte(username)))

		var counter int
		var err error

		if counterString == "" {
			counter = 0
		} else {
			counter, err = strconv.Atoi(counterString)

			if  err != nil {
				return errors.Wrap(err, "chicken counter is not a number")
			}
		}

		if counter == 0 {
			return errors.New("counter is already 0")
		}

		counter = counter - 1

		if counter == 0 {
			err = b.Delete([]byte(username))
		} else {
			counterString = strconv.Itoa(counter)
			err = b.Put([]byte(username), []byte(counterString))
		}

		if err == nil {
			userCounter = models.UserCounter{Username: username, Count: counter}
		}

		return err
	})

	return userCounter, err
}

// GetChickenStats Get chicken counters stats for everybody
func (s *Store) GetChickenStats(prefix string) ([]models.UserCounter, error) {
	bucketName := fmt.Sprintf("%s.%s", prefix, chickensBucket)

	userCounters := make([]models.UserCounter, 0)

	err := s.ensureBucketExists(bucketName)

	if err != nil {
		return userCounters, err
	}

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		c := b.Cursor()

		for key, value := c.First(); key != nil; key, value = c.Next() {
			counterString := string(value)
			counter, err := strconv.Atoi(counterString)

			if err != nil {
				return errors.Wrap(err, "chicken counter is not a number")
			}

			userCounter := models.UserCounter{Username: string(key), Count: counter}

			userCounters = append(userCounters, userCounter)
		}

		return nil
	})

	return userCounters, err
}

// Cakes

// AddCake Add to a cake counter for a username
func (s *Store) AddCake(username string, prefix string) (models.UserCounter, error) {
	bucketName := fmt.Sprintf("%s.%s", prefix, cakesBucket)

	var userCounter models.UserCounter

	err := s.ensureBucketExists(bucketName)

	if err != nil {
		return userCounter, err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		counterString := string(b.Get([]byte(username)))

		var counter int
		var err error

		if counterString == "" {
			counter = 0
		} else {
			counter, err = strconv.Atoi(counterString)

			if  err != nil {
				return errors.Wrap(err, "cake counter is not a number")
			}
		}

		counter = counter + 1
		counterString = strconv.Itoa(counter)

		err = b.Put([]byte(username), []byte(counterString))

		if err == nil {
			userCounter = models.UserCounter{Username: username, Count: counter}
		}

		return err
	})

	return userCounter, err
}

// FulfillCake Fulfill cake promise for a username
func (s *Store) FulfillCake(username string, prefix string) (models.UserCounter, error) {
	bucketName := fmt.Sprintf("%s.%s", prefix, cakesBucket)

	var userCounter models.UserCounter

	err := s.ensureBucketExists(bucketName)

	if err != nil {
		return userCounter, err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		counterString := string(b.Get([]byte(username)))

		var counter int
		var err error

		if counterString == "" {
			counter = 0
		} else {
			counter, err = strconv.Atoi(counterString)

			if  err != nil {
				return errors.Wrap(err, "cake counter is not a number")
			}
		}

		if counter == 0 {
			return errors.New("counter is already 0")
		}

		counter = counter - 1

		if counter == 0 {
			err = b.Delete([]byte(username))
		} else {
			counterString = strconv.Itoa(counter)
			err = b.Put([]byte(username), []byte(counterString))
		}

		if err == nil {
			userCounter = models.UserCounter{Username: username, Count: counter}
		}

		return err
	})

	return userCounter, err
}

// GetCakeStats Get cake counters stats for everybody
func (s *Store) GetCakeStats(prefix string) ([]models.UserCounter, error) {
	bucketName := fmt.Sprintf("%s.%s", prefix, cakesBucket)

	userCounters := make([]models.UserCounter, 0)

	err := s.ensureBucketExists(bucketName)

	if err != nil {
		return userCounters, err
	}

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		c := b.Cursor()

		for key, value := c.First(); key != nil; key, value = c.Next() {
			counterString := string(value)
			counter, err := strconv.Atoi(counterString)

			if err != nil {
				return errors.Wrap(err, "cake counter is not a number")
			}

			userCounter := models.UserCounter{Username: string(key), Count: counter}

			userCounters = append(userCounters, userCounter)
		}

		return nil
	})

	return userCounters, err
}

func (s *Store) ensureBucketExists(bucketName string) error {
	var err error

	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})

	return err
}