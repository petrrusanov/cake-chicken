package api

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/petrrusanov/cake-chicken/app/rest"
	"net/http"
	"strings"
)

func (s *Rest) addCake(w http.ResponseWriter, r *http.Request) {
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
		response := slackTextResponse{
			Text: fmt.Sprintf("Please provide a username"),
		}

		err = renderJSON(w, r, response)

		if err != nil {
			rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
			return
		}

		return
	}

	userID := matches[0][1]

	cakeCounter, err := s.DataStore.AddCake(userID, prefix)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "couldn't add cake")
		return
	}

	render.Status(r, http.StatusOK)

	var cakeText = strings.Repeat(":cake:", cakeCounter.Count)

	response := slackTextResponse{
		Text: fmt.Sprintf("Yay, more cakes are coming! <%s> have to bring %s", cakeCounter.Username, cakeText),
	}

	err = renderJSON(w, r, response)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
		return
	}
}

func (s *Rest) fulfillCake(w http.ResponseWriter, r *http.Request) {
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
		response := slackTextResponse{
			Text: fmt.Sprintf("Please provide a username"),
		}

		err = renderJSON(w, r, response)

		if err != nil {
			rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
			return
		}

		return
	}

	userID := matches[0][1]

	cakeCounter, err := s.DataStore.FulfillCake(userID, prefix)

	if err != nil {
		response := slackTextResponse{
			Text: fmt.Sprintf("Oops, for <%s> %s", userID, err.Error()),
		}

		err = renderJSON(w, r, response)

		if err != nil {
			rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
			return
		}

		return
	}

	render.Status(r, http.StatusOK)

	var cakeText string

	if cakeCounter.Count > 0 {
		cakeText = fmt.Sprintf("Still need to bring %s", strings.Repeat(":cake:", cakeCounter.Count))
	} else {
		cakeText = fmt.Sprintf("No more cakes to bring for <%s>", cakeCounter.Username)
	}

	response := slackTextResponse{
		Text: fmt.Sprintf("Thanks for the cake <%s>! %s", cakeCounter.Username, cakeText),
	}

	err = renderJSON(w, r, response)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
		return
	}
}
