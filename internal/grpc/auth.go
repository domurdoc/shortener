package grpc

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/model"
)

type ctxKey string

const userKey = ctxKey("user")

func CreateAuthIntercepter(a *auth.Auth) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		user, err := authenticataeOrRegisterAndLogin(ctx, a)
		if err != nil {
			var invalidTokenErr *auth.InvalidTokenError
			if errors.As(err, &invalidTokenErr) {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		return handler(attachUser(ctx, user), req)
	}
}

func authenticataeOrRegisterAndLogin(ctx context.Context, a *auth.Auth) (*model.User, error) {
	token, err := readToken(ctx)
	if err != nil {
		return nil, err
	}
	user, err := a.AuthenticateToken(ctx, token)
	if err != nil {
		var noTokenErr *auth.NoTokenError
		if errors.As(err, &noTokenErr) {
			user, err = a.Register(ctx)
			if err != nil {
				return nil, err
			}
			token, err := a.GenerateToken(ctx, user)
			if err != nil {
				return nil, err
			}
			if err := writeToken(ctx, token); err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}
	return user, nil
}

func readToken(ctx context.Context) (string, error) {
	var token string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		creds := md.Get("authorization")
		if len(creds) > 0 {
			token = creds[0]
		}
	}
	if token == "" {
		return "", fmt.Errorf("no token provided")
	}
	return token, nil
}

func writeToken(ctx context.Context, token string) error {
	md := metadata.Pairs("authorization", token)
	return grpc.SetHeader(ctx, md)
}

func attachUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func getUser(ctx context.Context) *model.User {
	return ctx.Value(userKey).(*model.User)
}
