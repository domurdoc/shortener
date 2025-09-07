package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/domurdoc/shortener/internal/httputils"
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
				contentType: httputils.ContentTypeTextPlain,
			},
		},
		{
			name:    "relative address",
			longURL: "/abcdef/",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: httputils.ContentTypeTextPlain,
			},
		},
		{
			name:    "not url-encoded",
			longURL: "/привет как дела?",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: httputils.ContentTypeTextPlain,
			},
		},
		{
			name:    "no-schema?",
			longURL: "yandex.com",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: httputils.ContentTypeTextPlain,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(service.New(repository.NewMem(), tt.baseURL))

			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.longURL))
			w := httptest.NewRecorder()
			handler.Shorten(w, r)

			resp := w.Result()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get(httputils.HeaderContentType))

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
			repo := repository.NewMem()
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

func TestShortener_ShortenJSON(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name        string
		body        string
		contentType string
		want        want
	}{
		{
			name:        "Invalid Content-Type: plain/text",
			contentType: httputils.ContentTypeTextPlain,
			want:        want{statusCode: http.StatusBadRequest},
		},
		{
			name:        "Invalid Content-Type: empty",
			contentType: "",
			want:        want{statusCode: http.StatusBadRequest},
		},
		{
			name:        "Invalid JSON: empty",
			body:        "",
			contentType: httputils.ContentTypeJSON,
			want:        want{statusCode: http.StatusBadRequest},
		},
		{
			name:        "Invalid JSON: broken",
			body:        `{"`,
			contentType: httputils.ContentTypeJSON,
			want:        want{statusCode: http.StatusBadRequest},
		},
		{
			name:        "Invalid JSON: null",
			body:        `null`,
			contentType: httputils.ContentTypeJSON,
			want:        want{statusCode: http.StatusBadRequest},
		},
		{
			name:        "Invalid JSON: no 'url' key",
			body:        `{"notaurl": "hello"}`,
			contentType: httputils.ContentTypeJSON,
			want:        want{statusCode: http.StatusBadRequest},
		},
		{
			name:        "Common request",
			body:        `{"url": "http://yandex.com"}`,
			contentType: httputils.ContentTypeJSON,
			want:        want{statusCode: http.StatusCreated},
		},
		{
			name:        "Extra JSON keys",
			body:        `{"url": "http://yandex.com", "extrakey": 123, "xxx": "18"}`,
			contentType: httputils.ContentTypeJSON,
			want:        want{statusCode: http.StatusCreated},
		},
		{
			name:        "'URL' key (uppercase)",
			body:        `{"url": "http://yandex.com"}`,
			contentType: httputils.ContentTypeJSON,
			want:        want{statusCode: http.StatusCreated},
		},
		{
			name:        "Ivalid url",
			body:        `{"url": "hello"}`,
			contentType: httputils.ContentTypeJSON,
			want:        want{statusCode: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New(service.New(repository.NewMem(), "http://localhost:8081"))

			r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.body))
			r.Header.Set(httputils.HeaderContentType, tt.contentType)
			w := httptest.NewRecorder()
			handler.ShortenJSON(w, r)

			resp := w.Result()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			if resp.StatusCode == http.StatusBadRequest {
				return
			}
			assert.Equal(t, httputils.ContentTypeJSON, resp.Header.Get(httputils.HeaderContentType))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			var respJSON Response
			err = json.Unmarshal(body, &respJSON)
			require.NoError(t, err)

			assert.NotEmpty(t, respJSON.Result)
		})
	}
}
