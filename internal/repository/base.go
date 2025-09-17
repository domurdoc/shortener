package repository

import "context"

type (
	Key   string
	Value string
)

type BatchItem struct {
	Key   Key
	Value Value
}

type Repo interface {
	Store(context.Context, Key, Value) error
	Fetch(context.Context, Key) (Value, error)
	StoreBatch(context.Context, []BatchItem) error
	Ping(context.Context) error
	Close() error
}

func SingleItemBatch(key Key, value Value) []BatchItem {
	return []BatchItem{{Key: key, Value: value}}
}
