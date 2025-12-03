package service

import (
	"database/sql"
	"time"

	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/generator"
	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/repository"
)

// Service encapsulates the business logic for managing shortened URLs.
// It handles operations such as shortening URLs, resolving short codes,
// batch processing, and asynchronous deletion of user records.
//
// The service uses a background worker pool to process deletions in batches,
// improving performance and reducing database load.
//
// Fields:
//   - baseURL: The base URL used to construct full short URLs (e.g., "https://short.ly").
//   - maxWorkers: The number of concurrent workers processing deletion batches.
//   - maxBatchSize: The maximum number of deletion records to process in a single batch.
//   - checkInterval: The time interval after which an incomplete batch is flushed and processed.
//   - deletedRecords: A channel used to queue user records marked for deletion.
//   - doneCh: A channel used to signal shutdown to all background goroutines.
//   - repo: The repository interface for data persistence operations.
//   - log: The logger instance for structured logging.
//   - db: Optional *sql.DB instance used for direct database health checks or transactions.
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
	gen            generator.Generator
}

func New(
	baseURL string,
	maxWorkers int,
	maxBatchSize int,
	checkInterval time.Duration,
	repo repository.RecordRepo,
	log *zap.SugaredLogger,
	db *sql.DB,
	gen generator.Generator,
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
		gen:            gen,
	}
	go d.serveDeletions()
	return d
}

func (s *Service) Close() error {
	close(s.doneCh)
	return nil
}
