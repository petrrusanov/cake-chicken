package api

import (
	"net/http"

	"github.com/boilerplate/backend/app/rest"
	"github.com/go-chi/render"

	"github.com/boilerplate/backend/app/utils"

	"github.com/boilerplate/backend/app/store/models"
)

func (s *Rest) createUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, hardBodyLimit), &user); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse user")
		return
	}

	token, err := utils.GenerateRandomString(64)

	if err != nil {
		http.Error(w, "failed to generate a token", http.StatusInternalServerError)
		return
	}

	token = utils.StrongHashValue(token, s.SharedSecret)

	user.Token = token

	newUser, err := s.DataStore.CreateUser(user)

	if err != nil {
		http.Error(w, "data store error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)

	err = renderJSON(w, r, newUser)

	if err != nil {
		http.Error(w, "json render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Rest) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := rest.MustGetUserInfo(r)
	err := renderJSON(w, r, user)

	if err != nil {
		http.Error(w, "json render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
