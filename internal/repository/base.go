package repository

type (
	Key   string
	Value string
)

type Repo interface {
	Store(Key, Value) error
	Fetch(Key) (Value, error)
}
