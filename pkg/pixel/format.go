package pixel

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Pixel struct {
	Site  uint64
	Event Event
}

type Event struct {
	Globals
	Name   string                 `json:"-"`
	Page   PageEvent              `json:"-"`
	Action ActionEvent            `json:"-"`
	More   map[string]interface{} `json:"-"`
}

type Globals struct {
	Timestamp time.Time `json:"-"`
	Visitor   string    `json:"-"`
}

type PageEvent struct {
	Name     string `json:"page"`
	Chapter1 string `json:"page_chapter1"`
	Chapter2 string `json:"page_chapter2"`
	Chapter3 string `json:"page_chapter3"`
}

type ActionEvent struct {
	Name     string `json:"action"`
	Type     string `json:"action_type"`
	Chapter1 string `json:"action_chapter1"`
	Chapter2 string `json:"action_chapter2"`
	Chapter3 string `json:"action_chapter3"`
}

func UnparseQuery(query []byte) (Pixel, error) {
	parts, err := url.Parse(string(query))

	if err != nil {
		return Pixel{}, errors.New("Invalid query string given")
	}

	if !parts.Query().Has("s") || !parts.Query().Has("p") {
		return Pixel{}, errors.New("Missing required fields `s` and/or `p`")
	}

	site, err := strconv.ParseUint(parts.Query().Get("s"), 10, 64)

	if err != nil {
		return Pixel{}, errors.New("Invalid `s` format")
	}

	var event *Event

	if err := json.Unmarshal([]byte(parts.Query().Get("p")), &event); err != nil {
		return Pixel{}, err
	}

	return Pixel{
		Site:  site,
		Event: *event,
	}, nil
}

func (pixel Pixel) String() string {
	return fmt.Sprintf("Pixel<site_id=%d, type=%s, visitor=%s, ts=%s>", pixel.Site, pixel.Event.Name, pixel.Event.Visitor, pixel.Event.Timestamp.UTC().Format(time.RFC3339))
}

func (event *Event) UnmarshalJSON(data []byte) error {
	var receiver map[string]interface{}

	if err := json.Unmarshal(data, &receiver); err != nil {
		return errors.New("Invalid `p` format")
	}

	if name, ok := receiver["event_name"].(string); ok {
		event.Name = name
	} else {
		return errors.New("Missing required field `event_name`")
	}

	fmt.Println(receiver)

	if visitor, ok := receiver["visitor"].(string); ok {
		event.Visitor = visitor
	} else {
		return errors.New("Missing required field `visitor`")
	}

	switch event.Name {
	case "":
		return errors.New("Invalid value for field `event_name`")
	case "page":
		if err := json.Unmarshal(data, &event.Page); err != nil {
			return errors.Join(errors.New("Invalid page event"), err)
		}
	case "action":
		if err := json.Unmarshal(data, &event.Action); err != nil {
			return errors.Join(errors.New("Invalid action event"), err)
		}
	}

	event.Timestamp = time.Now().UTC()

	if ts, ok := receiver["ts"].(string); ok {
		fmt.Println("ts", ts)
		if tsF, tsErr := strconv.ParseInt(ts, 10, 64); nil == tsErr {
			event.Timestamp = time.Unix(tsF, 0).UTC()
		}
	}

	event.More = make(map[string]interface{})

	var reservedKeys []string
	reservedKeys = append(reservedKeys, "event_name", "visitor", "ts")
	reservedEvents := make([]interface{}, 0)
	reservedEvents = append(reservedEvents, PageEvent{})
	reservedEvents = append(reservedEvents, ActionEvent{})

	for _, e := range reservedEvents {
		eventParsed, _ := json.Marshal(e)
		eventUnparsed := make(map[string]interface{}, 0)
		_ = json.Unmarshal(eventParsed, &eventUnparsed)
		for k := range eventUnparsed {
			reservedKeys = append(reservedKeys, k)
		}
	}

	for key, value := range receiver {
		isValid := true
		for _, invalid := range reservedKeys {
			if key == invalid {
				isValid = false
			}
		}

		if isValid {
			event.More[key] = value
		}
	}

	return nil
}
