package logger

import "go.uber.org/zap"

// TODO: do I need some kind of Log interface OR is it overkill??
func New(level string) *zap.SugaredLogger {
	lvl, _ := zap.ParseAtomicLevel(level)
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl := zap.Must(cfg.Build())
	return zl.Sugar()
}
