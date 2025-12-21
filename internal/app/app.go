package app

import (
	"errors"
	"fmt"
	_ "net/http/pprof"
	"time"

	"io/fs"

	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/auth/strategy"
	"github.com/domurdoc/shortener/internal/auth/transport"
	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/generator"
	"github.com/domurdoc/shortener/internal/logger"
	"github.com/domurdoc/shortener/internal/service"
	"github.com/domurdoc/shortener/internal/utils"
)

// App represents the main application structure that holds all components and configuration.
// It manages the lifecycle of repositories, services, authentication and audit logging.
type App struct {
	Config  *config.Config     // Options contains the application configuration.
	Log     *zap.SugaredLogger // Log is the logger instance used across the application.
	Service *service.Service   // Service is the core business logic service for URL operations.
	Repos   *Repositories      // Repositories manages database connections and repositories.
	Auth    *auth.Auth         // Auth manages authentication using JWT and cookies.
	Audit   *Audit             // Audit is the application-specific audit event handler.
	closer  utils.Closer
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
	a.closer.Register(
		// Sync() returns error when sink is console, so ignore it
		// https://github.com/uber-go/zap/issues/880
		func() error {
			err := a.Log.Sync()
			var pErr *fs.PathError
			if errors.As(err, &pErr) {
				return nil
			}
			return err
		},
	)
	return nil
}

func (a *App) initRepo() error {
	r, err := NewRepositories(a.Config, a.Log)
	if err != nil {
		return err
	}
	a.Repos = r
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
		a.Repos.Record,
		a.Log,
		a.Repos.DB,
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
	a.Auth = auth.New(strategy, transport, a.Repos.User)
	return nil
}

func (a *App) initAudit() error {
	a.Audit = NewAuditApp(a.Config, a.Log)
	a.closer.Register(a.Audit.Close)
	return nil
}
