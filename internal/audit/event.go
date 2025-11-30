package audit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/domurdoc/shortener/internal/model"
)

const actionShorten = "shorten"
const actionFollow = "follow"

type Event struct {
	TS     time.Time    `json:"-"`
	Action string       `json:"action"`
	UserID model.UserID `json:"-"`
	URL    string       `json:"url"`
}

func newEvent(action string, userID model.UserID, url string) *Event {
	return &Event{
		TS:     time.Now(),
		Action: action,
		UserID: userID,
		URL:    url,
	}
}

func (a *Audit) Shortened(userID model.UserID, url string) {
	e := newEvent(actionShorten, userID, url)
	a.Publish(e)
}
func (a *Audit) Followed(userID model.UserID, url string) {
	e := newEvent(actionFollow, userID, url)
	a.Publish(e)
}

func (e Event) MarshalJSON() ([]byte, error) {
	type EventAlias Event

	aliasValue := struct {
		EventAlias
		TS     int64  `json:"ts"`
		UserID string `json:"user_id,omitempty"`
	}{
		EventAlias: EventAlias(e),
		TS:         e.TS.Unix(),
		UserID:     fmt.Sprintf("%d", e.UserID),
	}
	return json.Marshal(aliasValue)
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type EventAlias Event

	aliasValue := &struct {
		*EventAlias
		TS     int64  `json:"ts"`
		UserID string `json:"user_id"`
	}{
		EventAlias: (*EventAlias)(e),
	}
	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}
	userID, err := strconv.Atoi(aliasValue.UserID)
	if err != nil {
		return err
	}
	e.UserID = model.UserID(userID)
	e.TS = time.Unix(aliasValue.TS, 0)
	return nil
}
