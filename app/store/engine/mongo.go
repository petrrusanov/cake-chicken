package engine

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-pkgz/mongo"
	"github.com/pkg/errors"
	"github.com/boilerplate/backend/app/store/models"
)

// Store extends engine interface with specific fields to the store
type Store struct {
	Interface

	connection *mongo.Connection
}

const (
	mongoUsers       = "users"
)

// NewMongo creates new store with specified params
func NewMongo(connection *mongo.Connection, bufferSize int, flushDuration time.Duration) (*Store, error) {
	result := Store{connection: connection}
	err := result.prepare()
	return &result, errors.Wrap(err, "failed to prepare mongo")
}

// prepare collections with all indexes
func (s *Store) prepare() error {
	return nil
}

// Close mongo store
func (s *Store) Close() error {
	return nil
}

func (s *Store) setLimitAndSkip(q *mgo.Query, limit, skip int) *mgo.Query {
	if limit <= 0 {
		limit = 1000
	}

	if skip < 0 {
		skip = 0
	}

	return q.Skip(skip).Limit(limit)
}

// FindUser find user by token
func (s *Store) FindUser(token string) (user models.User, err error) {
	err = s.connection.WithCustomCollection(mongoUsers, func(coll *mgo.Collection) error {
		query := bson.M{"token": token}
		return coll.Find(query).One(&user)
	})

	return user, err
}

// CreateUser create new user
func (s *Store) CreateUser(user models.User) (models.User, error) {
	err := s.connection.WithCustomCollection(mongoUsers, func(coll *mgo.Collection) error {
		return coll.Insert(&user)
	})

	return user, err
}
