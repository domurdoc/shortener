package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/domurdoc/shortener/internal/auth"

	"github.com/domurdoc/shortener/internal/auth/strategy"
	"github.com/domurdoc/shortener/internal/auth/transport"
	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/repository/mem"
	"github.com/domurdoc/shortener/internal/service"
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
			debugStrategy := strategy.NewDebug()
			bearerTransport := transport.NewBearer("Authorization")
			userRepo := mem.NewMemUserRepo()

			auth := auth.New(
				debugStrategy,
				bearerTransport,
				userRepo,
			)

			repo := mem.NewMemRecordRepo()
			handler := New(service.New(repo, nil, "", time.Second), auth, nil)

			if tt.want.statusCode == http.StatusTemporaryRedirect {
				user, _ := auth.Register(context.TODO())
				record := &model.BaseRecord{
					OriginalURL: model.OriginalURL(tt.want.location),
					ShortCode:   model.ShortCode(tt.shortCode),
				}
				err := repo.Store(context.TODO(), record, user.ID)
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
