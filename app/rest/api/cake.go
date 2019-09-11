package api

import (
	"errors"
	"fmt"
	"github.com/dimebox/cake-chicken/app/rest"
	"github.com/go-chi/render"
	"net/http"
)


/*
token=gIkuvaNzQIHg97ATvDxqgjtO
&team_id=T0001
&team_domain=example
&enterprise_id=E0001
&enterprise_name=Globular%20Construct%20Inc
&channel_id=C2147483705
&channel_name=test
&user_id=U2147483697
&user_name=Steve
&command=/weather
&text=94070
&response_url=https://hooks.slack.com/commands/1234/5678
&trigger_id=13345224609.738474920.8088930838d88f008e0
*/

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
		err = errors.New("username is missing")
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "text doesn't contain a username")
		return
	}

	userID := matches[0][1]

	cakeCounter, err := s.DataStore.AddCake(userID, prefix)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "couldn't add cake")
		return
	}

	render.Status(r, http.StatusOK)

	var cakeText string

	if cakeCounter.Count == 1 {
		cakeText = "cake"
	} else {
		cakeText = "cakes"
	}

	response := SlackTextResponse{
		Text: fmt.Sprintf("Yay, more cakes are coming! <%s> now owes %d %s", cakeCounter.Username, cakeCounter.Count, cakeText),
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
		err = errors.New("username is missing")
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "text doesn't contain a username")
		return
	}

	userID := matches[0][1]

	cakeCounter, err := s.DataStore.FulfillCake(userID, prefix)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "couldn't fulfill cake")
		return
	}

	render.Status(r, http.StatusOK)

	var cakeText string

	if cakeCounter.Count == 1 {
		cakeText = "cake"
	} else {
		cakeText = "cakes"
	}

	response := SlackTextResponse{
		Text: fmt.Sprintf("Hopefully it was tasty! <%s> now owes %d %s", cakeCounter.Username, cakeCounter.Count, cakeText),
	}

	err = renderJSON(w, r, response)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
		return
	}
}
