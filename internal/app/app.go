package app

import (
	"database/sql"
	"errors"
	"fmt"
	_ "net/http/pprof"
	"time"

	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/audit"
	"github.com/domurdoc/shortener/internal/audit/subscribers"
	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/auth/strategy"
	"github.com/domurdoc/shortener/internal/auth/transport"
	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/config/db"
	"github.com/domurdoc/shortener/internal/generator"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/domurdoc/shortener/internal/profiler"
	"github.com/domurdoc/shortener/internal/repository"
	dbRepo "github.com/domurdoc/shortener/internal/repository/db"
	fileRepo "github.com/domurdoc/shortener/internal/repository/file"
	"github.com/domurdoc/shortener/internal/repository/file/serializer"
	memRepo "github.com/domurdoc/shortener/internal/repository/mem"
	"github.com/domurdoc/shortener/internal/service"
	"github.com/domurdoc/shortener/internal/utils"
)

// App represents the main application structure that holds all components and configuration.
// It manages the lifecycle of repositories, services, authentication, audit logging, and profiling.
type App struct {
	Config         *config.Config                // Options contains the application configuration.
	RecordRepo     repository.RecordRepo         // RecordRepo is the repository for managing URL records.
	UserRepo       repository.UserRepo           // UserRepo is the repository for managing users.
	Log            *zap.SugaredLogger            // Log is the logger instance used across the application.
	Service        *service.Service              // Service is the core business logic service for URL operations.
	DB             *sql.DB                       // DB holds the database connection if using PostgreSQL.
	Auth           *auth.Auth                    // Auth manages authentication using JWT and cookies.
	Audit          *audit.Audit                  // Audit is the audit event dispatcher for logging actions.
	AuditFileSub   *subscribers.FileSubscriber   // AuditFileSub is the subscriber for writing audit logs to a file.
	AuditRemoteSub *subscribers.RemoteSubscriber // AuditRemoteSub is the subscriber for sending audit logs to a remote service.
	Profiler       *profiler.Profiler            // Profiler is the profiling server for monitoring performance.
	closer         utils.Closer
}

func New(cfg *config.Config) (*App, error) {
	a := &App{Config: cfg}
	if err := a.initLog(); err != nil {
		return nil, a.Close(err)
	}
	if err := a.initRepo(); err != nil {
		return nil, a.Close(err)
	}
	if err := a.initService(); err != nil {
		return nil, a.Close(err)
	}
	if err := a.initAuth(); err != nil {
		return nil, a.Close(err)
	}
	if err := a.initAudit(); err != nil {
		return nil, a.Close(err)
	}
	if err := a.initProfileServer(); err != nil {
		return nil, a.Close(err)
	}
	return a, nil
}

func (a *App) Close(reason error) error {
	return errors.Join(reason, a.closer.Close())
}

func (a *App) initLog() error {
	log, err := logger.New(a.Config.Logger.Level)
	if err != nil {
		return err
	}
	a.Log = log
	a.closer.Register(a.Log.Sync)
	return nil
}

func (a *App) initRepo() error {
	if a.Config.Repositories.DB.DSN != "" {
		pgDB, err := db.NewPG(a.Config.Repositories.DB.DSN)
		if err != nil {
			return err
		}
		a.DB = pgDB
		a.closer.Register(a.DB.Close)
		if err := db.MigratePG(pgDB); err != nil {
			return err
		}
		a.RecordRepo = dbRepo.NewDBRecordRepo(pgDB, db.NewPGArger)
		a.UserRepo = dbRepo.NewDBUserRepo(pgDB, db.NewPGArger)
	} else if a.Config.Repositories.File.Path != "" {
		jsonSerializer := serializer.NewJSONSerializer()
		repo, err := fileRepo.New(
			a.Config.Repositories.File.Path,
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
	var gen generator.Generator

	if a.Config.Generator.Constant.Value != "" {
		gen = generator.NewConstantGenerator(
			a.Config.Generator.Constant.Value,
		)
	} else if a.Config.Generator.Random.CharSet != "" && a.Config.Generator.Random.Length > 0 {
		gen = generator.NewRandomGenerator(
			a.Config.Generator.Random.CharSet,
			a.Config.Generator.Random.Length,
		)
	} else {
		return fmt.Errorf("failed to initialize Generator")
	}
	a.Service = service.New(
		a.Config.Service.BaseURL,
		a.Config.Service.DeleterMaxWorkers,
		a.Config.Service.DeleterMaxBatchSize,
		time.Duration(a.Config.Service.DeleterCheckInterval),
		a.RecordRepo,
		a.Log,
		a.DB,
		gen,
	)
	a.closer.Register(a.Service.Close)
	return nil
}

func (a *App) initAuth() error {

	strategy := strategy.NewJWT(
		a.Config.Auth.Strategy.JWTSecret,
		time.Duration(a.Config.Auth.Strategy.JWTDuration),
	)
	transport := transport.NewCookie(
		a.Config.Auth.Transport.CookieName,
		int(a.Config.Auth.Transport.CookieMaxAge.Seconds()),
		false,
	)
	a.Auth = auth.New(strategy, transport, a.UserRepo)
	return nil
}

func (a *App) initAudit() error {
	a.Audit = audit.New()

	if a.Config.Audit.File.Path != "" {
		a.AuditFileSub = subscribers.NewFile(
			a.Config.Audit.File.Path,
			a.Config.Audit.File.PoolSize,
			a.Config.Audit.File.MaxBatchSize,
			a.Config.Audit.File.BatchInterval,
			a.Log,
		)
		a.closer.Register(a.AuditFileSub.Close)
		a.Audit.Register(a.AuditFileSub)
	}
	if a.Config.Audit.Remote.URL != "" {
		a.AuditRemoteSub = subscribers.NewRemote(
			a.Config.Audit.Remote.URL,
			a.Log,
			a.Config.Audit.Remote.PoolSize,
		)
		a.closer.Register(a.AuditRemoteSub.Close)
		a.Audit.Register(a.AuditRemoteSub)
	}
	return nil
}

func (a *App) initProfileServer() error {
	if a.Config.Profiler.Address == "" {
		return nil
	}
	a.Profiler = profiler.New(a.Config.Profiler.Address)
	a.closer.Register(a.Profiler.Close)
	return nil
}
