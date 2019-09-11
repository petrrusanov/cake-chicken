package api

import (
	"errors"
	"fmt"
	"github.com/dimebox/cake-chicken/app/rest"
	"github.com/go-chi/render"
	"net/http"
)

func (s *Rest) addChicken(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse form")
		return
	}

	teamID := r.Form.Get("team_id")
	channelID := r.Form.Get("channel_id")

	if teamID == "" || channelID == "" {
		err = errors.New("invalid request")
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "team_id or channel_id are missing")
		return
	}

	prefix := fmt.Sprintf("%s.%s", teamID, channelID)
	text := r.Form.Get("text")

	matches := usernameRegexp.FindAllStringSubmatch(text,-1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		err = errors.New("username is missing")
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "text doesn't contain a username")
		return
	}

	userID := matches[0][1]

	chickenCounter, err := s.DataStore.AddChicken(userID, prefix)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "couldn't add chicken")
		return
	}

	render.Status(r, http.StatusCreated)

	err = renderJSON(w, r, chickenCounter)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
		return
	}
}

func (s *Rest) fulfillChicken(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse form")
		return
	}

	teamID := r.Form.Get("team_id")
	channelID := r.Form.Get("channel_id")

	if teamID == "" || channelID == "" {
		err = errors.New("invalid request")
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "team_id or channel_id are missing")
		return
	}

	prefix := fmt.Sprintf("%s.%s", teamID, channelID)
	text := r.Form.Get("text")

	matches := usernameRegexp.FindAllStringSubmatch(text,-1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		err = errors.New("username is missing")
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "text doesn't contain a username")
		return
	}

	userID := matches[0][1]

	chickenCounter, err := s.DataStore.FulfillChicken(userID, prefix)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "couldn't fulfill chicken")
		return
	}

	render.Status(r, http.StatusCreated)

	err = renderJSON(w, r, chickenCounter)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
		return
	}
}