package repository

import (
	"context"

	"github.com/domurdoc/shortener/internal/model"
)

// RecordRepo defines the interface for data access operations related to shortened URLs.
// It abstracts the storage, retrieval, and deletion of URL records and supports batch operations.
//
// Implementations of this interface are responsible for persisting and querying
// short code to original URL mappings, as well as handling user-associated records.
type RecordRepo interface {
	// Store saves a single URL record along with the associated user ID.
	// If the original URL already exists for the user, it may return model.OriginalURLExistsError.
	//
	// Parameters:
	//   - ctx: Context for timeout and cancellation control.
	//   - record: Pointer to the BaseRecord containing the original URL and short code.
	//   - userID: The ID of the user who owns the record.
	//
	// Returns:
	//   - nil on success.
	//   - An error if the operation fails, such as validation, constraint violation, or DB error.
	Store(context.Context, *model.BaseRecord, model.UserID) error

	// Fetch retrieves a URL record by its short code.
	//
	// Parameters:
	//   - ctx: Context for timeout and cancellation control.
	//   - shortCode: The short code string used to look up the record.
	//
	// Returns:
	//   - A pointer to the BaseRecord if found.
	//   - An error if the record does not exist or another issue occurs.
	Fetch(context.Context, model.ShortCode) (*model.BaseRecord, error)

	// FetchForUser retrieves all URL records associated with a given user.
	//
	// Parameters:
	//   - ctx: Context for timeout and cancellation control.
	//   - userID: The ID of the user whose records are being fetched.
	//
	// Returns:
	//   - A slice of BaseRecord instances owned by the user.
	//   - An error if the fetch operation fails.
	FetchForUser(context.Context, model.UserID) ([]model.BaseRecord, error)

	// StoreBatch saves multiple URL records in a single operation, associated with the given user.
	// It may partially succeed; if some records conflict, implementations may return BatchOriginalURLExistsError.
	//
	// Parameters:
	//   - ctx: Context for timeout and cancellation control.
	//   - records: A slice of BaseRecord to be stored.
	//   - userID: The ID of the user who owns the records.
	//
	// Returns:
	//   - nil on full success.
	//   - A model.BatchOriginalURLExistsError if some URLs already exist.
	//   - Other errors in case of system failure.
	StoreBatch(context.Context, []model.BaseRecord, model.UserID) error

	// Delete marks multiple user-record associations as deleted.
	// This is typically used for soft deletes or cleanup operations.
	//
	// Parameters:
	//   - ctx: Context for timeout and cancellation control.
	//   - userRecords: A slice of UserRecord, each specifying a user ID and short code to delete.
	//
	// Returns:
	//   - The number of records successfully marked as deleted.
	//   - An error if the operation fails partially or completely.
	Delete(context.Context, []model.UserRecord) (int, error)
}

// UserRepo defines the interface for user management operations.
// It supports user retrieval and creation, enabling user-centric features like ownership and authentication.
type UserRepo interface {
	// GetUser retrieves a user by their unique identifier.
	//
	// Parameters:
	//   - ctx: Context for timeout and cancellation control.
	//   - userID: The ID of the user to retrieve.
	//
	// Returns:
	//   - A pointer to the User if found.
	//   - An error if the user does not exist or another issue occurs.
	GetUser(context.Context, model.UserID) (*model.User, error)

	// CreateUser generates and stores a new user with a unique ID.
	// The implementation should ensure the generated ID is globally unique.
	//
	// Parameters:
	//   - ctx: Context for timeout and cancellation control.
	//
	// Returns:
	//   - A pointer to the newly created User.
	//   - An error if user creation fails.
	CreateUser(context.Context) (*model.User, error)
}
