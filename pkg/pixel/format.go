package pixel

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Pixel struct {
	Site  uint64
	Event Event
}

type Event struct {
	Name   string                 `json:"-"`
	Page   PageEvent              `json:"-"`
	Action ActionEvent            `json:"-"`
	More   map[string]interface{} `json:"-"`
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
	str := fmt.Sprintf("Pixel<site_id=%d,type:%s>", pixel.Site, pixel.Event.Name)

	switch pixel.Event.Name {
	case "page":
		data, _ := json.Marshal(pixel.Event.Page)
		str = fmt.Sprintf("%s%s", str, data)
	case "action":
		data, _ := json.Marshal(pixel.Event.Action)
		str = fmt.Sprintf("%s%s", str, data)
	}

	return str
}

func (event *Event) UnmarshalJSON(data []byte) error {
	var receiver map[string]interface{}

	if err := json.Unmarshal(data, &receiver); err != nil {
		return errors.New("Invalid `p` format")
	}

	switch receiver["event_name"] {
	case nil:
		return errors.New("Missing required event field `event_name`")
	case "":
		return errors.New("Missing required event field `event_name`")
	case "page":
		if err := json.Unmarshal(data, &event.Page); err != nil {
			return errors.Join(errors.New("Invalid page event"), err)
		}
	case "action":
		if err := json.Unmarshal(data, &event.Action); err != nil {
			return errors.Join(errors.New("Invalid action event"), err)
		}
	}
	event.Name = receiver["event_name"].(string)
	event.More = make(map[string]interface{})

	var reservedKeys []string
	reservedKeys = append(reservedKeys, "event_name")
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
