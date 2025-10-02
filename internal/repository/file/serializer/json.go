package serializer

import (
	"encoding/json"

	"github.com/domurdoc/shortener/internal/model"
)

type jsonRecord struct {
	ShortURL    model.ShortCode   `json:"short_url"`
	OriginalURL model.OriginalURL `json:"original_url"`
}

type jsonOwnership struct {
	UserID    model.UserID    `json:"user_id"`
	ShortCode model.ShortCode `json:"short_url"`
}

type jsonSnapshot struct {
	Records   []jsonRecord    `json:"records"`
	Ownership []jsonOwnership `json:"ownership"`
}

func toJSONRecord(r model.Record) jsonRecord {
	return jsonRecord{
		ShortURL:    r.ShortCode,
		OriginalURL: r.OriginalURL,
	}
}

func fromJSONRecord(jr jsonRecord) model.Record {
	return model.Record{
		ShortCode:   jr.ShortURL,
		OriginalURL: jr.OriginalURL,
	}
}

func toJSONSnapshot(r *Snapshot) jsonSnapshot {
	jsonRecords := make([]jsonRecord, 0, len(r.Records))
	for _, r := range r.Records {
		jr := toJSONRecord(r)
		jsonRecords = append(jsonRecords, jr)
	}

	jsonOwnerships := make([]jsonOwnership, 0, len(r.Ownership))
	for _, o := range r.Ownership {
		jo := jsonOwnership(o)
		jsonOwnerships = append(jsonOwnerships, jo)
	}

	return jsonSnapshot{
		Records:   jsonRecords,
		Ownership: jsonOwnerships,
	}
}

func fromJSONSnapshot(js *jsonSnapshot) *Snapshot {
	records := make([]model.Record, 0, len(js.Records))
	for _, jr := range js.Records {
		r := fromJSONRecord(jr)
		records = append(records, r)
	}

	ownership := make([]Ownership, 0, len(js.Ownership))
	for _, jo := range js.Ownership {
		o := Ownership(jo)
		ownership = append(ownership, o)
	}

	return &Snapshot{
		Records:   records,
		Ownership: ownership,
	}
}

type JSONSerializer struct{}

func (s *JSONSerializer) Dump(snapshot *Snapshot) ([]byte, error) {
	jsonSnapshot := toJSONSnapshot(snapshot)
	return json.Marshal(jsonSnapshot)
}

func (s *JSONSerializer) Load(data []byte) (*Snapshot, error) {
	if len(data) == 0 {
		return nil, nil
	}
	var jsonSnapshot jsonSnapshot
	if err := json.Unmarshal(data, &jsonSnapshot); err != nil {
		return nil, err
	}
	return fromJSONSnapshot(&jsonSnapshot), nil
}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}
