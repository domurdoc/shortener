package mem

import (
	"context"
	"maps"
	"sync"

	"slices"

	"github.com/domurdoc/shortener/internal/model"
)

// MemUserRepo is an in-memory implementation of the UserRepo interface.
// It stores users in a map and uses a mutex for concurrent access safety.
type MemUserRepo struct {
	storage map[model.UserID]model.User
	mu      sync.Mutex
}

// NewMemUserRepo creates and returns a new instance of MemUserRepo with an empty user map.
func NewMemUserRepo() *MemUserRepo {
	return &MemUserRepo{storage: make(map[model.UserID]model.User)}
}

// GetUser retrieves a user by their unique ID from the in-memory storage.
// It is safe for concurrent use.
//
// Parameters:
//   - ctx: Context for potential future use (e.g., timeouts).
//   - userID: The ID of the user to retrieve.
//
// Returns:
//   - A pointer to the User if found.
//   - A *model.UserNotFoundError if the user does not exist.
func (m *MemUserRepo) GetUser(ctx context.Context, userID model.UserID) (*model.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	user, ok := m.storage[userID]
	if !ok {
		return nil, &model.UserNotFoundError{UserID: userID}
	}
	return &user, nil
}

// CreateUser generates a new user with a unique ID and stores it in memory.
// The ID is assigned sequentially based on the current highest ID.
// It is safe for concurrent use.
//
// Parameters:
//   - ctx: Context for potential future use.
//
// Returns:
//   - A pointer to the newly created User.
//   - nil error on success.
func (m *MemUserRepo) CreateUser(ctx context.Context) (*model.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	nextUserID := model.UserID(1)
	if len(m.storage) > 0 {
		nextUserID = slices.Max(slices.Collect(maps.Keys(m.storage))) + 1
	}

	user := model.User{ID: nextUserID}
	m.storage[nextUserID] = user
	return &user, nil
}

// CountUsers returns the total number of users in the repository.
func (m *MemUserRepo) CountUsers(ctx context.Context) (int, error) {
	return len(m.storage), nil
}
