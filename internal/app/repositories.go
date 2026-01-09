package app

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/config/db"
	"github.com/domurdoc/shortener/internal/repository"
	dbRepo "github.com/domurdoc/shortener/internal/repository/db"
	fileRepo "github.com/domurdoc/shortener/internal/repository/file"
	"github.com/domurdoc/shortener/internal/repository/file/serializer"
	memRepo "github.com/domurdoc/shortener/internal/repository/mem"
	"github.com/domurdoc/shortener/internal/utils"
)

type Repositories struct {
	Record repository.RecordRepo
	User   repository.UserRepo
	DB     *sql.DB
	closer utils.Closer
}

func NewRepositories(cfg *config.Config, log *zap.SugaredLogger) (*Repositories, error) {
	r := &Repositories{}

	if cfg.RepositoryDSN != "" {
		pgDB, err := db.NewPG(cfg.RepositoryDSN)
		if err != nil {
			return nil, err
		}
		r.DB = pgDB
		r.closer.Register(r.DB.Close)
		if err := db.MigratePG(pgDB); err != nil {
			return nil, err
		}
		r.Record = dbRepo.NewDBRecordRepo(pgDB, db.NewPGArger)
		r.User = dbRepo.NewDBUserRepo(pgDB, db.NewPGArger)
	} else if cfg.RepositoryFilePath != "" {
		jsonSerializer := serializer.NewJSONSerializer()
		repo, err := fileRepo.New(
			cfg.RepositoryFilePath,
			jsonSerializer,
		)
		if err != nil {
			return nil, err
		}
		r.Record = repo
		r.User = memRepo.NewMemUserRepo()
	} else {
		r.Record = memRepo.NewMemRecordRepo()
		r.User = memRepo.NewMemUserRepo()
	}
	return r, nil
}

func (r *Repositories) Close() error {
	return r.closer.Close()
}
