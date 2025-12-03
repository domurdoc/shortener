package db

// Arger is an interface representing an argument provider for database operations.
// It is used to sequentially retrieve string arguments, typically for query construction or scanning.
// The implementing type manages internal state to return the next argument on each call to Next.
//
// Example usage might include parsing command-line arguments, iterating over query parameters,
// or reading from a configuration source.
type Arger interface {
	// Next returns the next argument as a string.
	// It returns an empty string when no more arguments are available.
	// The caller is responsible for knowing how many arguments to expect.
	Next() string
}
