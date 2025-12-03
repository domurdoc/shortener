package model

// OriginalURL represents the full original URL string that has been shortened.
// It is a typed alias for string to enforce type safety in the domain model.
type OriginalURL string

// ShortCode represents the unique key used to identify a shortened URL.
// It is a typed alias for string (e.g., "abc123") and is used in redirects.
type ShortCode string

// ShortURL represents the full shortened URL as a string (e.g., "https://short.ly/abc123").
// It is a typed alias for string used for returning complete URLs to clients.
type ShortURL string

// BaseRecord holds the core data of a shortened URL entry.
// It maps a short code to its original URL.
type BaseRecord struct {
	// ShortCode is the unique identifier for the shortened URL.
	ShortCode ShortCode
	// OriginalURL is the full URL that the short code points to.
	OriginalURL OriginalURL
}

// UserRecord represents the association between a user and a short code.
// It is used to track ownership and enable per-user operations like deletion.
type UserRecord struct {
	// ShortCode is the short code owned by the user.
	ShortCode ShortCode
	// UserID is the identifier of the user who owns the short code.
	UserID UserID
}

// URLRecord represents a user-facing shortened URL entry.
// It contains both the full short URL and the original URL.
type URLRecord struct {
	// ShortURL is the complete shortened URL (base URL + short code).
	ShortURL ShortURL
	// OriginalURL is the full original URL that the short URL points to.
	OriginalURL OriginalURL
}
