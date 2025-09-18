package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/domurdoc/shortener/internal/repository"
	"github.com/domurdoc/shortener/internal/repository/mem"
	"github.com/domurdoc/shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			repo := mem.New()
			handler := New(service.New(repo, ""))

			if tt.want.statusCode == http.StatusTemporaryRedirect {
				err := repo.Store(context.TODO(), repository.Key(tt.shortCode), repository.Value(tt.want.location))
				require.NoError(t, err)
			}

			r := httptest.NewRequest(http.MethodGet, "/{shortCode}", nil)
			w := httptest.NewRecorder()
			r.SetPathValue("shortCode", tt.shortCode)

			handler.Retrieve(w, r)

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
