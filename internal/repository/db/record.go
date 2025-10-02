package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/domurdoc/shortener/internal/config/db"
	"github.com/domurdoc/shortener/internal/model"
)

type DBRecordRepo struct {
	db       *sql.DB
	newArger func() db.Arger
}

func NewDBRecordRepo(db *sql.DB, newArger func() db.Arger) *DBRecordRepo {
	return &DBRecordRepo{db, newArger}
}

const (
	queryInsertRecord = `
INSERT INTO records (key, value) VALUES (%s, %s)
ON CONFLICT (value) DO UPDATE SET key = records.key
RETURNING id, key
`
	queryInsertOwnership = `
INSERT INTO ownership (user_id, record_id) VALUES (%s, %s)
ON CONFLICT (user_id, record_id) DO NOTHING
`
	queryFetchRecord = `
SELECT value FROM records WHERE key = %s
`
	queryFetchForUser = `
SELECT key, value FROM records r JOIN ownership o ON r.id = o.record_id
WHERE o.user_id = %s
`
)

func (r *DBRecordRepo) Store(ctx context.Context, record *model.Record, userID model.UserID) error {
	var arger db.Arger

	arger = r.newArger()
	insertRecordQuery := fmt.Sprintf(queryInsertRecord, arger.Next(), arger.Next())
	arger = r.newArger()
	insertOwnershipQuery := fmt.Sprintf(queryInsertOwnership, arger.Next(), arger.Next())

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(
		ctx,
		insertRecordQuery,
		record.ShortCode,
		record.OriginalURL,
	)

	var recordID int
	var shortCode model.ShortCode

	err = row.Scan(
		&recordID,
		&shortCode,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		insertOwnershipQuery,
		userID,
		recordID,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	if shortCode != record.ShortCode {
		return &model.OriginalURLExistsError{
			OriginalURL: record.OriginalURL,
			ShortCode:   shortCode,
		}
	}
	return nil
}

func (r *DBRecordRepo) StoreBatch(ctx context.Context, records []model.Record, userID model.UserID) error {
	var arger db.Arger

	arger = r.newArger()
	insertRecordQuery := fmt.Sprintf(queryInsertRecord, arger.Next(), arger.Next())
	arger = r.newArger()
	insertOwnershipQuery := fmt.Sprintf(queryInsertOwnership, arger.Next(), arger.Next())

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertRecordStmt, err := tx.PrepareContext(ctx, insertRecordQuery)
	if err != nil {
		return err
	}
	defer insertRecordStmt.Close()

	insertOwnershipStmt, err := tx.PrepareContext(ctx, insertOwnershipQuery)
	if err != nil {
		return err
	}
	defer insertOwnershipStmt.Close()

	var batchError model.BatchError
	for pos, record := range records {
		row := insertRecordStmt.QueryRowContext(
			ctx,
			record.ShortCode,
			record.OriginalURL,
		)
		var recordID int
		var shortCode model.ShortCode
		err = row.Scan(
			&recordID,
			&shortCode,
		)
		if err != nil {
			return err
		}
		_, err = insertOwnershipStmt.ExecContext(ctx,
			userID,
			recordID,
		)
		if err != nil {
			return err
		}
		if shortCode != record.ShortCode {
			valueErr := &model.OriginalURLExistsError{
				OriginalURL: record.OriginalURL,
				ShortCode:   shortCode,
				BatchPos:    pos,
			}
			batchError = append(batchError, valueErr)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	if len(batchError) != 0 {
		return batchError
	}
	return nil
}

func (r *DBRecordRepo) Fetch(ctx context.Context, shortCode model.ShortCode) (*model.Record, error) {
	record := model.Record{ShortCode: shortCode}

	arger := r.newArger()
	query := fmt.Sprintf(queryFetchRecord, arger.Next())

	row := r.db.QueryRowContext(
		ctx,
		query,
		shortCode,
	)
	err := row.Scan(&record.OriginalURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &model.ShortCodeNotFoundError{ShortCode: shortCode}
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *DBRecordRepo) FetchForUser(ctx context.Context, userID model.UserID) ([]model.Record, error) {
	var records []model.Record

	arger := r.newArger()
	fetchForUserQuery := fmt.Sprintf(queryFetchForUser, arger.Next())

	rows, err := r.db.QueryContext(ctx, fetchForUserQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		record := model.Record{}
		if err := rows.Scan(
			&record.ShortCode,
			&record.OriginalURL,
		); err != nil {
			return records, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return records, err
	}
	return records, nil
}
