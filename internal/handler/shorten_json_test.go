package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
			contentType: httputil.ContentTypeTextPlain,
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
			contentType: httputil.ContentTypeJSON,
			want:        want{statusCode: http.StatusBadRequest},
		},
		{
			name:        "Invalid JSON: broken",
			body:        `{"`,
			contentType: httputil.ContentTypeJSON,
			want:        want{statusCode: http.StatusBadRequest},
		},
		{
			name:        "Invalid JSON: null",
			body:        `null`,
			contentType: httputil.ContentTypeJSON,
			want:        want{statusCode: http.StatusBadRequest},
		},
		{
			name:        "Invalid JSON: no 'url' key",
			body:        `{"notaurl": "hello"}`,
			contentType: httputil.ContentTypeJSON,
			want:        want{statusCode: http.StatusBadRequest},
		},
		{
			name:        "Common request",
			body:        `{"url": "http://yandex.com"}`,
			contentType: httputil.ContentTypeJSON,
			want:        want{statusCode: http.StatusCreated},
		},
		{
			name:        "Extra JSON keys",
			body:        `{"url": "http://yandex.com", "extrakey": 123, "xxx": "18"}`,
			contentType: httputil.ContentTypeJSON,
			want:        want{statusCode: http.StatusCreated},
		},
		{
			name:        "'URL' key (uppercase)",
			body:        `{"url": "http://yandex.com"}`,
			contentType: httputil.ContentTypeJSON,
			want:        want{statusCode: http.StatusCreated},
		},
		{
			name:        "Ivalid url",
			body:        `{"url": "hello"}`,
			contentType: httputil.ContentTypeJSON,
			want:        want{statusCode: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugStrategy := strategy.NewDebug()
			bearerTransport := transport.NewBearer("Authorization")
			userRepo := mem.NewMemUserRepo()

			a := auth.New(
				debugStrategy,
				bearerTransport,
				userRepo,
			)
			service := service.New(
				"http://localhost:8081",
				1,
				1,
				time.Second,
				mem.NewMemRecordRepo(),
				nil,
				nil,
			)
			handler := New(service)

			r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.body))
			r.Header.Set(httputil.HeaderContentType, tt.contentType)
			w := httptest.NewRecorder()

			user, err := a.AuthenticateOrRegisterAndLogin(context.TODO(), w, r)
			assert.NoError(t, err)

			handler.ShortenJSON(w, auth.AttachUser(r, user))

			resp := w.Result()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			if resp.StatusCode == http.StatusBadRequest {
				return
			}
			assert.Equal(t, httputil.ContentTypeJSON, resp.Header.Get(httputil.HeaderContentType))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			err = resp.Body.Close()
			require.NoError(t, err)

			var respJSON jsonResponse
			err = json.Unmarshal(body, &respJSON)
			require.NoError(t, err)

			assert.NotEmpty(t, respJSON.Result)
		})
	}
}
