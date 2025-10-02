package mem

import (
	"context"
	"maps"
	"sync"

	"slices"

	"github.com/domurdoc/shortener/internal/model"
)

type MemUserRepo struct {
	storage map[model.UserID]model.User
	mu      sync.Mutex
}

func NewMemUserRepo() *MemUserRepo {
	return &MemUserRepo{storage: make(map[model.UserID]model.User)}
}

func (m *MemUserRepo) GetUser(ctx context.Context, userID model.UserID) (*model.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	user, ok := m.storage[userID]
	if !ok {
		return nil, &model.UserNotFoundError{UserID: userID}
	}
	return &user, nil
}

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
