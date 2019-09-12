package api

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/petrrusanov/cake-chicken/app/rest"
	"net/http"
	"strings"
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
		response := slackTextResponse{
			Text: fmt.Sprintf("Please provide a username"),
			ResponseType: EphemeralResponse,
		}

		err = renderJSON(w, r, response)

		if err != nil {
			rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
			return
		}

		return
	}

	userID := matches[0][1]

	chickenCounter, err := s.DataStore.AddChicken(userID, prefix)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "couldn't add chicken")
		return
	}

	render.Status(r, http.StatusOK)

	chickenText := strings.Repeat(":poultry_leg:", chickenCounter.Count)

	response := slackTextResponse{
		Text: fmt.Sprintf("Chicken is on <%s>! %s", chickenCounter.Username, chickenText),
		ResponseType: InChannelResponse,
	}

	err = renderJSON(w, r, response)

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
		response := slackTextResponse{
			Text: fmt.Sprintf("Please provide a username"),
			ResponseType: EphemeralResponse,
		}

		err = renderJSON(w, r, response)

		if err != nil {
			rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
			return
		}

		return
	}

	userID := matches[0][1]

	chickenCounter, err := s.DataStore.FulfillChicken(userID, prefix)

	if err != nil {
		response := slackTextResponse{
			Text: fmt.Sprintf("Oops, for <%s> %s", userID, err.Error()),
			ResponseType: EphemeralResponse,
		}

		err = renderJSON(w, r, response)

		if err != nil {
			rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
			return
		}

		return
	}

	render.Status(r, http.StatusOK)

	var chickenText string

	if chickenCounter.Count > 0 {
		chickenText = fmt.Sprintf("Still need to bring %s", strings.Repeat(":poultry_leg:", chickenCounter.Count))
	} else {
		chickenText = fmt.Sprintf("No more chicken to bring for <%s>", chickenCounter.Username)
	}

	response := slackTextResponse{
		Text: fmt.Sprintf("Thanks for the chicken <%s>! %s", chickenCounter.Username, chickenText),
		ResponseType: InChannelResponse,
	}

	err = renderJSON(w, r, response)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
		return
	}
}