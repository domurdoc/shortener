package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/repository/mem"
	"github.com/domurdoc/shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			handler := New(service.New(mem.New(), "http://localhost:8081"))

			r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.body))
			r.Header.Set(httputil.HeaderContentType, tt.contentType)
			w := httptest.NewRecorder()
			handler.ShortenJSON(w, r)

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
