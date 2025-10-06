package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
SELECT value, NOT EXISTS(SELECT 1 FROM ownership o WHERE o.record_id = r.id) AS is_deleted FROM records r WHERE key = %s
`
	queryFetchForUser = `
SELECT key, value FROM records r JOIN ownership o ON r.id = o.record_id
WHERE o.user_id = %s
`
	queryDeleteOwnership = `
DELETE FROM
	ownership
WHERE
	(user_id, record_id) IN (
		SELECT
			o.user_id,
			o.record_id
		FROM
			ownership o
		JOIN
			records r
		ON
			r.id = o.record_id
		WHERE
			(o.user_id, r.key) IN (%s)
	)
`
)

func (r *DBRecordRepo) Store(ctx context.Context, record *model.BaseRecord, userID model.UserID) error {
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

func (r *DBRecordRepo) StoreBatch(ctx context.Context, records []model.BaseRecord, userID model.UserID) error {
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

	var batchError model.BatchOriginalURLExistsError
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

func (r *DBRecordRepo) Fetch(ctx context.Context, shortCode model.ShortCode) (*model.BaseRecord, error) {
	record := model.BaseRecord{ShortCode: shortCode}
	var isDeleted bool

	arger := r.newArger()
	query := fmt.Sprintf(queryFetchRecord, arger.Next())

	row := r.db.QueryRowContext(
		ctx,
		query,
		shortCode,
	)

	err := row.Scan(
		&record.OriginalURL,
		&isDeleted,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &model.ShortCodeNotFoundError{ShortCode: shortCode}
	}
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, &model.ShortCodeDeletedError{ShortCode: shortCode}
	}
	return &record, nil
}

func (r *DBRecordRepo) FetchForUser(ctx context.Context, userID model.UserID) ([]model.BaseRecord, error) {
	var records []model.BaseRecord

	arger := r.newArger()
	fetchForUserQuery := fmt.Sprintf(queryFetchForUser, arger.Next())

	rows, err := r.db.QueryContext(ctx, fetchForUserQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		record := model.BaseRecord{}
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

func (r *DBRecordRepo) Delete(ctx context.Context, records []model.UserRecord) (int, error) {
	arger := r.newArger()

	values := make([]string, 0, len(records))
	args := make([]any, 0, len(records))

	for _, record := range records {
		values = append(values, fmt.Sprintf("(%s, %s)", arger.Next(), arger.Next()))
		args = append(args, record.UserID, record.ShortCode)
	}

	deleteOwnershipQuery := fmt.Sprintf(queryDeleteOwnership, strings.Join(values, ","))
	res, err := r.db.ExecContext(ctx, deleteOwnershipQuery, args...)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	return int(count), err
}
