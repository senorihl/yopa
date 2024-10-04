package pixel

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnparseQuery(t *testing.T) {
	var tests = []struct {
		name  string
		input []byte
		want  Pixel
		err   string
	}{
		{"Should throw err if invalid string", []byte("Totally\r\nInvalid"), Pixel{}, "Invalid query string given"},
		{"Should throw err if empty", []byte(""), Pixel{}, "Missing required fields `s` and/or `p`"},
		{"Should throw err if invalid keys", []byte("test=false"), Pixel{}, "Missing required fields `s` and/or `p`"},
		{"Should throw err if missing `p`", []byte("?s=1337"), Pixel{}, "Missing required fields `s` and/or `p`"},
		{"Should throw err if missing `s`", []byte("?p={\"test\":\"test\"}"), Pixel{}, "Missing required fields `s` and/or `p`"},
		{"Should throw err if invalid `s` format",
			[]byte("?s=-1337&p={\"test\":\"test\"}"),
			Pixel{},
			"Invalid `s` format"},
		{"Should throw err if invalid `p` format",
			[]byte("?s=1337&p=invalid_json_string"),
			Pixel{},
			"invalid character 'i' looking for beginning of value"},
		{"Should throw err if missing `p.event_name`",
			[]byte("?s=1337&p={\"test\":\"test\"}"),
			Pixel{},
			"Missing required field `event_name`"},
		{"Should throw err if missing `p.visitor`",
			[]byte("?s=1337&p={\"event_name\":\"test\"}"),
			Pixel{},
			"Missing required field `visitor`"},
		{"Should throw err if empty `p.event_name`",
			[]byte("?s=1337&p={\"test\":\"\",\"visitor\":\"visitor\"}"),
			Pixel{},
			"Missing required field `event_name`"},
		{"Should handle custom event as map",
			[]byte("?s=1337&p={\"event_name\":\"test\",\"visitor\":\"visitor\",\"ts\":\"1727740800\"}"),
			Pixel{Site: 1337, Event: Event{Name: "test", Globals: Globals{Timestamp: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC).UTC(), Visitor: "visitor"}, Page: PageEvent{}, More: make(map[string]interface{})}},
			""},
		{"Should handle page event correctly",
			[]byte("?s=1337&p={\"event_name\":\"page\",\"page\":\"page\",\"visitor\":\"visitor\",\"ts\":\"1727740800\"}"),
			Pixel{Site: 1337, Event: Event{Name: "page", Globals: Globals{Timestamp: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC).UTC(), Visitor: "visitor"}, Page: PageEvent{Name: "page"}, More: make(map[string]interface{})}},
			""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := UnparseQuery(tt.input)
			assert.Equal(t, tt.want, res)

			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
