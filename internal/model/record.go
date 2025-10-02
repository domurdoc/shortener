package model

type (
	OriginalURL string
	ShortCode   string
	ShortURL    string
)

type Record struct {
	ShortCode   ShortCode
	OriginalURL OriginalURL
}

type URLRecord struct {
	ShortURL    ShortURL
	OriginalURL OriginalURL
}
