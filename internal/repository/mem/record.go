package mem

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/domurdoc/shortener/internal/model"
)

// MemRecordRepo is an in-memory implementation of the RecordRepo interface.
// It stores URL records and user associations using multiple maps for efficient lookups
// by short code, user ID, and original URL. It is thread-safe using a mutex.
type MemRecordRepo struct {
	// ShortCodeRecords maps short codes to their base records.
	ShortCodeRecords map[model.ShortCode]model.BaseRecord
	// ShortCodeUserIDS maps short codes to a nested map of user IDs and their associated records.
	ShortCodeUserIDS map[model.ShortCode]map[model.UserID]model.BaseRecord
	// UserIDRecords maps user IDs to a nested map of short codes and their associated records.
	UserIDRecords map[model.UserID]map[model.ShortCode]model.BaseRecord
	// OriginalURLRecords maps original URLs to their base records for deduplication.
	OriginalURLRecords map[model.OriginalURL]model.BaseRecord
	// mu protects all maps from concurrent access.
	mu sync.Mutex
}

// NewMemRecordRepo creates and returns a new instance of MemRecordRepo
// with all internal maps initialized.
func NewMemRecordRepo() *MemRecordRepo {
	return &MemRecordRepo{
		ShortCodeRecords:   make(map[model.ShortCode]model.BaseRecord),
		ShortCodeUserIDS:   make(map[model.ShortCode]map[model.UserID]model.BaseRecord),
		UserIDRecords:      make(map[model.UserID]map[model.ShortCode]model.BaseRecord),
		OriginalURLRecords: make(map[model.OriginalURL]model.BaseRecord),
	}
}

// Store saves a single URL record associated with the given user.
// If the original URL already exists, it returns an *model.OriginalURLExistsError.
// It delegates to StoreBatch for implementation.
//
// Parameters:
//   - ctx: Context for potential future use.
//   - record: The BaseRecord to store.
//   - userID: The ID of the user creating the record.
//
// Returns:
//   - nil on success.
//   - *model.OriginalURLExistsError if the original URL already exists.
//   - Other errors if storage fails.
func (r *MemRecordRepo) Store(ctx context.Context, record *model.BaseRecord, userID model.UserID) error {
	err := r.StoreBatch(ctx, []model.BaseRecord{*record}, userID)
	var batchURLExistsErr model.BatchOriginalURLExistsError
	if errors.As(err, &batchURLExistsErr) {
		return batchURLExistsErr[0]
	}
	return err
}

// StoreBatch saves multiple URL records associated with the given user.
// It checks for existing short codes and original URLs to prevent duplication.
// If any URL already exists, it returns a BatchOriginalURLExistsError listing conflicts.
//
// Parameters:
//   - ctx: Context for potential future use.
//   - records: Slice of BaseRecord to store.
//   - userID: The ID of the user creating the records.
//
// Returns:
//   - nil if all records are stored successfully.
//   - model.BatchOriginalURLExistsError if some URLs already exist.
//   - Error if a short code collision occurs.
func (r *MemRecordRepo) StoreBatch(ctx context.Context, records []model.BaseRecord, userID model.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, record := range records {
		if _, exists := r.ShortCodeRecords[record.ShortCode]; exists {
			return fmt.Errorf("ShortCode %s already exists", record.ShortCode)
		}
	}
	var batchURLExistsErr model.BatchOriginalURLExistsError
	for pos, record := range records {
		existingRecord, exists := r.OriginalURLRecords[record.OriginalURL]
		if !exists {
			r.OriginalURLRecords[record.OriginalURL] = record
			r.ShortCodeRecords[record.ShortCode] = record
			r.ShortCodeUserIDS[record.ShortCode] = make(map[model.UserID]model.BaseRecord)
		} else if existingRecord.ShortCode != record.ShortCode {
			record.ShortCode = existingRecord.ShortCode
			urlExistsErr := &model.OriginalURLExistsError{
				OriginalURL: record.OriginalURL,
				ShortCode:   existingRecord.ShortCode,
				BatchPos:    pos,
			}
			batchURLExistsErr = append(batchURLExistsErr, urlExistsErr)
		}
		if _, ok := r.UserIDRecords[userID]; !ok {
			r.UserIDRecords[userID] = make(map[model.ShortCode]model.BaseRecord)
		}
		r.UserIDRecords[userID][record.ShortCode] = record
		r.ShortCodeUserIDS[record.ShortCode][userID] = record
	}
	if len(batchURLExistsErr) != 0 {
		return batchURLExistsErr
	}
	return nil
}

// Fetch retrieves a BaseRecord by its short code.
// Returns an error if the short code does not exist or has been deleted (no associated users).
//
// Parameters:
//   - ctx: Context for potential future use.
//   - shortCode: The short code to look up.
//
// Returns:
//   - Pointer to the BaseRecord if found and active.
//   - *model.ShortCodeNotFoundError if not found.
//   - *model.ShortCodeDeletedError if the code exists but was deleted.
func (r *MemRecordRepo) Fetch(ctx context.Context, shortCode model.ShortCode) (*model.BaseRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, exists := r.ShortCodeRecords[shortCode]
	if !exists {
		return nil, &model.ShortCodeNotFoundError{ShortCode: shortCode}
	}
	userIDS := r.ShortCodeUserIDS[shortCode]
	if len(userIDS) == 0 {
		return nil, &model.ShortCodeDeletedError{ShortCode: shortCode}
	}
	return &record, nil
}

// FetchForUser retrieves all active URL records created by the specified user.
//
// Parameters:
//   - ctx: Context for potential future use.
//   - userID: The ID of the user whose records are requested.
//
// Returns:
//   - A slice of BaseRecord instances owned by the user.
//   - nil if the user has no records.
func (r *MemRecordRepo) FetchForUser(ctx context.Context, userID model.UserID) ([]model.BaseRecord, error) {
	originalURLRecords, ok := r.UserIDRecords[userID]
	if !ok {
		return nil, nil
	}
	return slices.Collect(maps.Values(originalURLRecords)), nil
}

// Delete removes user associations from records (soft delete).
// It removes the user from the short code's user list and the user's record list.
//
// Parameters:
//   - ctx: Context for potential future use.
//   - records: Slice of UserRecord indicating which user-short code pairs to delete.
//
// Returns:
//   - The number of associations successfully removed.
//   - nil error (errors are ignored for partial failures).
func (r *MemRecordRepo) Delete(ctx context.Context, records []model.UserRecord) (int, error) {
	counter := 0
	for _, record := range records {
		userIDS, ok := r.ShortCodeUserIDS[record.ShortCode]
		if !ok {
			continue
		}
		counter++
		delete(userIDS, record.UserID)
		shortCodeRecords, ok := r.UserIDRecords[record.UserID]
		if !ok {
			continue
		}
		delete(shortCodeRecords, record.ShortCode)
	}
	return counter, nil
}
