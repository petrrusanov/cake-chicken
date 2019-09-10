package engine

import "github.com/boilerplate/backend/app/store/models"

// Interface combines all store interfaces
type Interface interface {
	User
}

// User user profile methods
type User interface {
	FindUser(token string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
}
