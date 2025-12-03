package generator

// Generator defines the interface for generating short identifiers.
// It abstracts the logic for creating unique, compact strings used as short URLs.
// Implementations may use various strategies such as random strings, hashing, or sequence-based IDs.
type Generator interface {
	// Generate produces a new short ID string.
	// Returns the generated ID and nil error on success.
	// Returns an empty string and error if generation fails (e.g., collision, entropy exhaustion).
	Generate() (string, error)
}
