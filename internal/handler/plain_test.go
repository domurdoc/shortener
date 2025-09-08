package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/repository"
	"github.com/domurdoc/shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortener_Shorten(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name    string
		longURL string
		baseURL string
		want    want
	}{
		{
			name:    "generic",
			longURL: "http://yandex.com/abcdef/",
			baseURL: "http://localhost:8080",
			want: want{
				statusCode:  http.StatusCreated,
				contentType: httputil.ContentTypeTextPlain,
			},
		},
		{
			name:    "relative address",
			longURL: "/abcdef/",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: httputil.ContentTypeTextPlain,
			},
		},
		{
			name:    "not url-encoded",
			longURL: "/привет как дела?",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: httputil.ContentTypeTextPlain,
			},
		},
		{
			name:    "no-schema?",
			longURL: "yandex.com",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: httputil.ContentTypeTextPlain,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(service.New(repository.NewMemRepo(), tt.baseURL))

			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.longURL))
			w := httptest.NewRecorder()
			handler.Shorten(w, r)

			resp := w.Result()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get(httputil.HeaderContentType))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			if tt.want.statusCode == http.StatusCreated {
				parsedURL, err := url.Parse(string(body))
				assert.NoError(t, err)
				assert.Contains(t, parsedURL.String(), tt.baseURL)
				assert.NotEqual(t, parsedURL.Path, "")
			}
		})
	}
}

func TestShortener_Retrieve(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name      string
		shortCode string
		want      want
	}{
		{
			name:      "not found",
			shortCode: "x",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:      "found",
			shortCode: "x",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "http://yandex.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewMemRepo()
			handler := New(service.New(repo, ""))

			if tt.want.statusCode == http.StatusTemporaryRedirect {
				err := repo.Store(repository.Key(tt.shortCode), repository.Value(tt.want.location))
				require.NoError(t, err)
			}

			r := httptest.NewRequest(http.MethodGet, "/{shortCode}", nil)
			w := httptest.NewRecorder()
			r.SetPathValue("shortCode", tt.shortCode)

			handler.GetByShortCode(w, r)

			resp := w.Result()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			err := resp.Body.Close()
			require.NoError(t, err)

			if tt.want.statusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
			}
		})
	}
}
