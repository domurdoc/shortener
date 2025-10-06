package serializer

import "github.com/domurdoc/shortener/internal/model"

type Ownership struct {
	UserID    model.UserID
	ShortCode model.ShortCode
}

type Snapshot struct {
	Records   []model.BaseRecord
	Ownership []Ownership
}

type Serializer interface {
	Dump(snapshot *Snapshot) ([]byte, error)
	Load([]byte) (*Snapshot, error)
}
