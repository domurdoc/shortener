package model

// UserID uniquely identifies a user.
type UserID int

// User represents a system user with a unique ID.
type User struct {
	ID UserID
}
