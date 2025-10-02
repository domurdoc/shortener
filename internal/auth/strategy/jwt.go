package strategy

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/domurdoc/shortener/internal/model"
	"github.com/domurdoc/shortener/internal/repository"
)

type JWTStrategy struct {
	secretKey string
	tokenExp  time.Duration
}

type claims struct {
	jwt.RegisteredClaims
	UserID model.UserID
}

func NewJWT(secretKey string, tokenExp time.Duration) *JWTStrategy {
	return &JWTStrategy{secretKey: secretKey, tokenExp: tokenExp}
}

func (s *JWTStrategy) WriteToken(ctx context.Context, user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExp)),
		},
		UserID: user.ID,
	})
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JWTStrategy) ReadToken(ctx context.Context, tokenString string, repo repository.UserRepo) (*model.User, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			return []byte(s.secretKey), nil
		})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return repo.GetUser(ctx, claims.UserID)
}
