package service

import (
	"database/sql"
	"time"

	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/repository"
)

type Service struct {
	baseURL        string
	maxWorkers     int
	maxBatchSize   int
	checkInterval  time.Duration
	deletedRecords chan model.UserRecord
	doneCh         chan struct{}
	repo           repository.RecordRepo
	log            *zap.SugaredLogger
	db             *sql.DB
}

func New(
	baseURL string,
	maxWorkers int,
	maxBatchSize int,
	checkInterval time.Duration,
	repo repository.RecordRepo,
	log *zap.SugaredLogger,
	db *sql.DB,
) *Service {
	d := &Service{
		baseURL:        baseURL,
		maxWorkers:     maxWorkers,
		maxBatchSize:   maxBatchSize,
		checkInterval:  checkInterval,
		deletedRecords: make(chan model.UserRecord),
		doneCh:         make(chan struct{}),
		repo:           repo,
		log:            log,
		db:             db,
	}
	go d.serveDeletions()
	return d
}

func (s *Service) Close() error {
	close(s.doneCh)
	return nil
}
