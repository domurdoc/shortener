package file

import (
	"encoding/json"
	"strconv"

	"github.com/domurdoc/shortener/internal/repository"
)

type jsonRecord struct {
	ID          int              `json:"UUID"`
	ShortURL    repository.Key   `json:"short_url"`
	OriginalURL repository.Value `json:"original_url"`
}

func (r jsonRecord) MarshalJSON() ([]byte, error) {
	type jsonRecordAlias jsonRecord

	aliasValue := struct {
		jsonRecordAlias
		ID string `json:"UUID"`
	}{
		jsonRecordAlias: jsonRecordAlias(r),
		ID:              strconv.Itoa(r.ID),
	}
	return json.Marshal(aliasValue)
}

func (r *jsonRecord) UnmarshalJSON(data []byte) (err error) {
	type jsonRecordAlias jsonRecord

	aliasValue := &struct {
		*jsonRecordAlias
		ID string `json:"UUID"`
	}{
		jsonRecordAlias: (*jsonRecordAlias)(r),
	}
	if err = json.Unmarshal(data, aliasValue); err != nil {
		return err
	}
	r.ID, err = strconv.Atoi(aliasValue.ID)
	return
}

func toJSONRecord(r record) jsonRecord {
	return jsonRecord{
		ID:          r.ID,
		ShortURL:    r.Key,
		OriginalURL: r.Value,
	}
}

func fromJSONRecord(jr jsonRecord) record {
	return record{
		ID:    jr.ID,
		Key:   jr.ShortURL,
		Value: jr.OriginalURL,
	}
}

type JSONSerializer struct{}

func (s *JSONSerializer) Dump(records []record) ([]byte, error) {
	jsonRecords := make([]jsonRecord, 0, len(records))
	for _, r := range records {
		jsonRecords = append(jsonRecords, toJSONRecord(r))
	}
	return json.Marshal(jsonRecords)
}

func (s *JSONSerializer) Load(data []byte) ([]record, error) {
	var jsonRecords []jsonRecord
	var records []record
	if len(data) == 0 {
		return records, nil
	}
	if err := json.Unmarshal(data, &jsonRecords); err != nil {
		return nil, err
	}
	records = make([]record, 0, len(jsonRecords))
	for _, jr := range jsonRecords {
		records = append(records, fromJSONRecord(jr))
	}
	return records, nil
}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}
