package models

import (
	"encoding/json"
)

type Route struct {
	OriginalUrl string
	ShortPath   string
	VisitCount  int
}

func (r Route) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Route) UnMarshalBinary(text []byte) error {
	return json.Unmarshal(text, &r)
}
