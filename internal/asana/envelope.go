package asana

import "encoding/json"

type responseEnvelope struct {
	Data     json.RawMessage `json:"data"`
	NextPage *Page           `json:"next_page"`
}
