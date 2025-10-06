package model

type (
	OriginalURL string
	ShortCode   string
	ShortURL    string
)

type BaseRecord struct {
	ShortCode   ShortCode
	OriginalURL OriginalURL
}

type UserRecord struct {
	ShortCode ShortCode
	UserID    UserID
}

type URLRecord struct {
	ShortURL    ShortURL
	OriginalURL OriginalURL
}
