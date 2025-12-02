package app

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/audit"
	"github.com/domurdoc/shortener/internal/audit/subscribers"
	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/auth/strategy"
	"github.com/domurdoc/shortener/internal/auth/transport"
	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/config/db"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/domurdoc/shortener/internal/repository"
	dbRepo "github.com/domurdoc/shortener/internal/repository/db"
	fileRepo "github.com/domurdoc/shortener/internal/repository/file"
	"github.com/domurdoc/shortener/internal/repository/file/serializer"
	memRepo "github.com/domurdoc/shortener/internal/repository/mem"
	"github.com/domurdoc/shortener/internal/service"
)

type App struct {
	Options         *config.Options
	RecordRepo      repository.RecordRepo
	UserRepo        repository.UserRepo
	Log             *zap.SugaredLogger
	Service         *service.Service
	DB              *sql.DB
	Auth            *auth.Auth
	Audit           *audit.Audit
	AuditFileSub    *subscribers.FileSubscriber
	AuditRemoteSub  *subscribers.RemoteSubscriber
	ProfileServer   *http.Server
	profileServerWG sync.WaitGroup
}

func New() (*App, error) {
	a := &App{Options: config.LoadOptions()}
	if err := a.initRepo(); err != nil {
		return nil, errors.Join(err, a.Close())
	}
	if err := a.initLog(); err != nil {
		return nil, errors.Join(err, a.Close())
	}
	if err := a.initService(); err != nil {
		return nil, errors.Join(err, a.Close())
	}
	if err := a.initAuth(); err != nil {
		return nil, errors.Join(err, a.Close())
	}
	if err := a.initAudit(); err != nil {
		return nil, errors.Join(err, a.Close())
	}
	if err := a.initProfileServer(); err != nil {
		return nil, errors.Join(err, a.Close())
	}
	return a, nil
}

func (a *App) Close() error {
	var errs []error

	if a.ProfileServer != nil {
		errs = append(errs, a.ProfileServer.Shutdown(context.Background()))
		a.profileServerWG.Wait()
	}
	if a.AuditFileSub != nil {
		errs = append(errs, a.AuditFileSub.Close())
	}
	if a.AuditRemoteSub != nil {
		errs = append(errs, a.AuditRemoteSub.Close())
	}
	if a.Service != nil {
		errs = append(errs, a.Service.Close())
	}
	if a.Log != nil {
		errs = append(errs, a.Log.Sync())
	}
	if a.DB != nil {
		errs = append(errs, a.DB.Close())
	}
	return errors.Join(errs...)
}

func (a *App) initLog() error {
	log, err := logger.New(a.Options.LogLevel.String())
	if err != nil {
		return err
	}
	a.Log = log
	return nil
}

func (a *App) initRepo() error {
	if a.Options.DatabaseDSN.String() != "" {
		pgDB, err := db.NewPG(a.Options.DatabaseDSN.String())
		if err != nil {
			return err
		}
		a.DB = pgDB
		if err := db.MigratePG(pgDB); err != nil {
			return err
		}
		a.RecordRepo = dbRepo.NewDBRecordRepo(pgDB, db.NewPGArger)
		a.UserRepo = dbRepo.NewDBUserRepo(pgDB, db.NewPGArger)
	} else if a.Options.FileStoragePath.String() != "" {
		jsonSerializer := serializer.NewJSONSerializer()
		repo, err := fileRepo.New(
			a.Options.FileStoragePath.String(),
			jsonSerializer,
		)
		if err != nil {
			return err
		}
		a.RecordRepo = repo
		a.UserRepo = memRepo.NewMemUserRepo()
	} else {
		a.RecordRepo = memRepo.NewMemRecordRepo()
		a.UserRepo = memRepo.NewMemUserRepo()
	}
	return nil
}

func (a *App) initService() error {
	a.Service = service.New(
		a.Options.BaseURL.String(),
		int(a.Options.DeleterMaxWorkers),
		int(a.Options.DeleterMaxBatchSize),
		time.Duration(a.Options.DeleterCheckInterval),
		a.RecordRepo,
		a.Log,
		a.DB,
	)
	return nil
}

func (a *App) initAuth() error {
	strategy := strategy.NewJWT(
		a.Options.JWTSecret.String(),
		time.Duration(a.Options.JWTDuration),
	)
	transport := transport.NewCookie(
		a.Options.CookieName.String(),
		int(time.Duration(a.Options.CookieMaxAge).Seconds()),
		false,
	)
	a.Auth = auth.New(strategy, transport, a.UserRepo)
	return nil
}

func (a *App) initAudit() error {
	a.Audit = audit.New()

	if a.Options.AuditFile.String() != "" {
		a.AuditFileSub = subscribers.NewFile(
			a.Options.AuditFile.String(),
			int(a.Options.AuditFilePoolSize),
			int(a.Options.AuditFileMaxBatchSize),
			time.Duration(a.Options.AuditFileBatchInterval),
			a.Log,
		)
		a.Audit.Register(a.AuditFileSub)
	}
	if a.Options.AuditURL.String() != "" {
		a.AuditRemoteSub = subscribers.NewRemote(
			a.Options.AuditURL.String(),
			a.Log,
			int(a.Options.AuditRemotePoolSize),
		)
		a.Audit.Register(a.AuditRemoteSub)
	}
	return nil
}

func (a *App) initProfileServer() error {
	if a.Options.ProfileAddr.String() == "" {
		return nil
	}
	// https://stackoverflow.com/a/42533360
	server := &http.Server{Addr: a.Options.ProfileAddr.String()}

	a.profileServerWG.Add(1)
	go func() {
		defer a.profileServerWG.Done()

		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("ProfileServer.ListenAndServe(): %v", err)
		}
	}()
	a.ProfileServer = server
	return nil
}
