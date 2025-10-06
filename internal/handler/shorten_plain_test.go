package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/auth/strategy"
	"github.com/domurdoc/shortener/internal/auth/transport"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/repository/mem"
	"github.com/domurdoc/shortener/internal/service"
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
			debugStrategy := strategy.NewDebug()
			bearerTransport := transport.NewBearer("Authorization")
			userRepo := mem.NewMemUserRepo()

			auth := auth.New(
				debugStrategy,
				bearerTransport,
				userRepo,
			)
			service := service.New(
				tt.baseURL,
				1,
				1,
				time.Second,
				mem.NewMemRecordRepo(),
				nil,
				nil,
			)
			handler := New(service, auth)

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
