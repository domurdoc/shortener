package app

import (
	"go.uber.org/zap"

	"github.com/domurdoc/shortener/internal/audit"
	"github.com/domurdoc/shortener/internal/audit/subscribers"
	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/utils"
)

type Audit struct {
	Audit     *audit.Audit
	FileSub   *subscribers.FileSubscriber
	RemoteSub *subscribers.RemoteSubscriber
	closer    utils.Closer
}

func NewAuditApp(cfg *config.Config, log *zap.SugaredLogger) *Audit {
	a := &Audit{
		Audit: audit.New(),
	}
	if cfg.Audit.File.Path != "" {
		a.FileSub = subscribers.NewFile(
			cfg.Audit.File.Path,
			cfg.Audit.File.PoolSize,
			cfg.Audit.File.MaxBatchSize,
			cfg.Audit.File.BatchInterval,
			log,
		)
		a.closer.Register(a.FileSub.Close)
		a.Audit.Register(a.FileSub)
	}
	if cfg.Audit.Remote.URL != "" {
		a.RemoteSub = subscribers.NewRemote(
			cfg.Audit.Remote.URL,
			log,
			cfg.Audit.Remote.PoolSize,
		)
		a.closer.Register(a.RemoteSub.Close)
		a.Audit.Register(a.RemoteSub)
	}
	return a
}

func (a *Audit) Close() error {
	return a.closer.Close()
}
