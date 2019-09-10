package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/boilerplate/backend/app/store/models"
)

type contextKey string

// MustGetUserInfo fails if can't extract user data from the request.
// should be called from authed controllers only
func MustGetUserInfo(r *http.Request) models.User {
	user, err := GetUserInfo(r)
	if err != nil {
		panic(err)
	}
	return user
}

// GetUserInfo returns user from request context
func GetUserInfo(r *http.Request) (user models.User, err error) {
	ctx := r.Context()

	if ctx == nil {
		return models.User{}, errors.New("no info about user")
	}

	ctxUser := ctx.Value(contextKey("user"))

	if u, ok := ctxUser.(models.User); ok {
		return u, nil
	}

	return models.User{}, errors.New("user can't be parsed")
}

// SetUserInfo sets user into request context
func SetUserInfo(r *http.Request, user models.User) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKey("user"), user)
	return r.WithContext(ctx)
}
