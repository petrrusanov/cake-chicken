package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/boilerplate/backend/app/store/engine"
	"github.com/boilerplate/backend/app/store/models"
)

// DataStore wraps engine.Interface with additional methods
type DataStore struct {
	engine.Interface
}

// CreateUser prepares user and forward to Interface.CreateUser
func (s *DataStore) CreateUser(user models.User) (preparedUser models.User, err error) {
	if preparedUser, err = s.prepareNewUser(user); err != nil {
		return preparedUser, errors.Wrap(err, "failed to prepare user")
	}

	return s.Interface.CreateUser(preparedUser)
}

// prepareNewUser sets new user fields, hashing and sanitizing data
func (s *DataStore) prepareNewUser(user models.User) (models.User, error) {
	// fill ID and time if empty
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	return user, nil
}
