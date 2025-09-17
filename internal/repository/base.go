package repository

import "context"

type (
	Key   string
	Value string
)

type Repo interface {
	Store(context.Context, Key, Value) error
	Fetch(context.Context, Key) (Value, error)
	Ping(context.Context) error
	Close() error
}
