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

	stats := map[string]string{}

	for _, counter := range cakeStats {
		cakeText := fmt.Sprintf("%s", strings.Repeat(":cake:", counter.Count))
		stats[counter.Username] = cakeText
	}

	for _, counter := range chickenStats {
		chickenText := fmt.Sprintf("%s", strings.Repeat(":poultry_leg:", counter.Count))

		userStats := stats[counter.Username]

		if userStats != "" {
			userStats = fmt.Sprintf("%s %s", userStats, chickenText)
		} else {
			userStats = chickenText
		}

		stats[counter.Username] = userStats
	}

	var builder strings.Builder

	builder.WriteString("*Stats:*\n")

	if len(stats) == 0 {
		builder.WriteString("There is no cake or chicken to bring :cry:")
	} else {
		for username, text := range stats {
			builder.WriteString(fmt.Sprintf("<%s> %s\n", username, text))
		}
	}

	response := slackTextResponse{
		Text: builder.String(),
	}

	err = renderJSON(w, r, response)

	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "json render error")
		return
	}
}