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
	if cfg.AuditFilePath != "" {
		a.FileSub = subscribers.NewFile(
			cfg.AuditFilePath,
			cfg.AuditFilePoolSize,
			cfg.AuditFileMaxBatchSize,
			cfg.AuditFileBatchInterval,
			log,
		)
		a.closer.Register(a.FileSub.Close)
		a.Audit.Register(a.FileSub)
	}
	if cfg.AuditRemoteURL != "" {
		a.RemoteSub = subscribers.NewRemote(
			cfg.AuditRemoteURL,
			log,
			cfg.AuditRemotePoolSize,
		)
		a.closer.Register(a.RemoteSub.Close)
		a.Audit.Register(a.RemoteSub)
	}
	return a
}

func (a *Audit) Close() error {
	return a.closer.Close()
}
