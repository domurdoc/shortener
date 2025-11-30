package config

type Options struct {
	Addr                   NetAddress
	BaseURL                URL
	LogLevel               LogLevel
	FileStoragePath        String
	DatabaseDSN            String
	JWTSecret              String
	JWTDuration            Duration
	CookieName             String
	CookieMaxAge           Duration
	DeleterMaxWorkers      Integer
	DeleterMaxBatchSize    Integer
	DeleterCheckInterval   Duration
	AuditFile              String
	AuditFilePoolSize      Integer
	AuditFileMaxBatchSize  Integer
	AuditFileBatchInterval Duration
	AuditURL               String
	AuditRemotePoolSize    Integer
}

func New(
	addr,
	baseURL,
	logLevel,
	storagePath,
	databaseDSN,
	jwtSecret,
	jwtDuration,
	cookieName,
	cookieMaxAge,
	deleterMaxWorkers,
	deleterMaxBatchSize,
	deleterCheckInterval,
	auditFile,
	auditFilePoolSize,
	auditFileMaxBatchSize,
	auditFileBatchInterval,
	auditURL,
	auditRemotePoolSize string,
) *Options {
	options := Options{}
	setOptionFromString(&options.BaseURL, baseURL)
	setOptionFromString(&options.Addr, addr)
	setOptionFromString(&options.LogLevel, logLevel)
	setOptionFromString(&options.FileStoragePath, storagePath)
	setOptionFromString(&options.DatabaseDSN, databaseDSN)
	setOptionFromString(&options.JWTSecret, jwtSecret)
	setOptionFromString(&options.JWTDuration, jwtDuration)
	setOptionFromString(&options.CookieMaxAge, cookieMaxAge)
	setOptionFromString(&options.CookieName, cookieName)
	setOptionFromString(&options.DeleterMaxWorkers, deleterMaxWorkers)
	setOptionFromString(&options.DeleterMaxBatchSize, deleterMaxBatchSize)
	setOptionFromString(&options.DeleterCheckInterval, deleterCheckInterval)
	setOptionFromString(&options.AuditFile, auditFile)
	setOptionFromString(&options.AuditFilePoolSize, auditFilePoolSize)
	setOptionFromString(&options.AuditFileMaxBatchSize, auditFileMaxBatchSize)
	setOptionFromString(&options.AuditFileBatchInterval, auditFileBatchInterval)
	setOptionFromString(&options.AuditURL, auditURL)
	setOptionFromString(&options.AuditRemotePoolSize, auditRemotePoolSize)
	return &options
}

type option interface {
	Set(string) error
}

func setOptionFromString(o option, value string) {
	if err := o.Set(value); err != nil {
		panic(err)
	}
}
