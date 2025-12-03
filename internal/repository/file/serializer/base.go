package serializer

import "github.com/domurdoc/shortener/internal/model"

// Ownership represents the association between a user and a short code.
// It is used in the Snapshot to track which users own which shortened URLs.
type Ownership struct {
	// UserID is the identifier of the user who owns the short code.
	UserID model.UserID
	// ShortCode is the identifier of the shortened URL.
	ShortCode model.ShortCode
}

// Snapshot represents the complete state of the URL storage system.
// It contains all URL records and their ownership mappings for persistence.
type Snapshot struct {
	// Records holds all stored base records indexed by short code.
	Records []model.BaseRecord
	// Ownership lists all user-to-short-code associations.
	Ownership []Ownership
}

// Serializer defines the interface for converting a Snapshot to and from raw bytes.
// This allows different serialization formats (e.g., JSON, Gob) to be used interchangeably.
type Serializer interface {
	// Dump serializes the given Snapshot into a byte slice.
	//
	// Parameters:
	//   - snapshot: The Snapshot to serialize.
	//
	// Returns:
	//   - The serialized data as a byte slice.
	//   - An error if serialization fails.
	Dump(snapshot *Snapshot) ([]byte, error)

	// Load deserializes a byte slice into a Snapshot.
	//
	// Parameters:
	//   - data: The raw serialized data to deserialize.
	//
	// Returns:
	//   - A pointer to the reconstructed Snapshot, or nil if data is empty.
	//   - An error if deserialization fails.
	Load([]byte) (*Snapshot, error)
}
