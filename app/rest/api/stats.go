package api

import (
	"errors"
	"fmt"
	"github.com/petrrusanov/cake-chicken/app/rest"
	"net/http"
	"strings"
)

func (s *Rest) stats(w http.ResponseWriter, r *http.Request) {
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

	cakeStats, err := s.DataStore.GetCakeStats(prefix)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get cake stats")
		return
	}

	chickenStats, err := s.DataStore.GetChickenStats(prefix)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get chicken stats")
		return
	}

	var builder strings.Builder

	builder.WriteString("*Cakes:*")

	if len(cakeStats) == 0 {
		builder.WriteString(" none\n")
	} else {
		builder.WriteString("\n")
	}

	for _, counter := range cakeStats {
		var cakeText string

		if counter.Count == 1 {
			cakeText = "cake"
		} else {
			cakeText = "cakes"
		}

		_, err := fmt.Fprintf(&builder, "<%s> owes %d %s\n", counter.Username, counter.Count, cakeText)

		if err != nil {
			rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't format cake stats")
		}
	}

	builder.WriteString("\n*Chickens:*")

	if len(chickenStats) == 0 {
		builder.WriteString(" none\n")
	} else {
		builder.WriteString("\n")
	}

	for _, counter := range chickenStats {
		var chickenText string

		if counter.Count == 1 {
			chickenText = "chicken"
		} else {
			chickenText = "chickens"
		}

		_, err := fmt.Fprintf(&builder, "<%s> owes %d %s\n", counter.Username, counter.Count, chickenText)

		if err != nil {
			rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't format chicken stats")
		}
	}

	response := SlackTextResponse{
		Text: builder.String(),
	}

	err = renderJSON(w, r, response)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
		return
	}
}