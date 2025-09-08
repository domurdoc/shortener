package repository

type (
	Key   string
	Value string
)

type Repo interface {
	Store(Key, Value) error
	Fetch(Key) (Value, error)
}

type Record struct {
	ID    int
	Key   Key
	Value Value
}

type Serializer interface {
	Dump([]Record) ([]byte, error)
	Load([]byte) ([]Record, error)
}
