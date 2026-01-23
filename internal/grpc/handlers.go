package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/domurdoc/shortener/api"
	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/model"
)

type ShortenerServer struct {
	api.UnimplementedShortenerServiceServer
	app *app.App
}

func NewShortenerServiceServer(a *app.App) *ShortenerServer {
	return &ShortenerServer{app: a}
}

func (s *ShortenerServer) ShortenURL(ctx context.Context, in *api.URLShortenRequest) (*api.URLShortenResponse, error) {
	var response api.URLShortenResponse

	user := getUser(ctx)

	longURL := in.GetUrl()
	shortURL, err := s.app.Service.Shorten(ctx, user, longURL)
	s.app.Audit.Audit.Shortened(user.ID, longURL)
	var invalidURLErr *model.InvalidURLError
	if errors.As(err, &invalidURLErr) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	var urlExistsErr *model.OriginalURLExistsError
	if err != nil && !errors.As(err, &urlExistsErr) {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response.SetResult(shortURL)
	if err != nil {
		return &response, status.Error(codes.AlreadyExists, err.Error())
	}
	return &response, nil
}

func (s *ShortenerServer) ExpandURL(ctx context.Context, in *api.URLExpandRequest) (*api.URLExpandResponse, error) {
	shortCode := in.GetId()
	longURL, err := s.app.Service.GetByShortCode(ctx, shortCode)
	s.app.Audit.Audit.Followed(model.UserID(0), longURL)
	var notFoundErr *model.ShortCodeNotFoundError
	if errors.As(err, &notFoundErr) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	var isDeletedErr *model.ShortCodeDeletedError
	if errors.As(err, &isDeletedErr) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return api.URLExpandResponse_builder{Result: &longURL}.Build(), nil
}

func (s *ShortenerServer) ListUserURLs(ctx context.Context, e *emptypb.Empty) (*api.UserURLsResponse, error) {
	var response api.UserURLsResponse

	user := getUser(ctx)

	urlRecords, err := s.app.Service.GetForUser(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	urlDatas := make([]*api.URLData, 0, len(urlRecords))
	for _, r := range urlRecords {
		urlDatas = append(urlDatas, api.URLData_builder{
			ShortUrl:    proto.String(string(r.ShortURL)),
			OriginalUrl: proto.String(string(r.OriginalURL)),
		}.Build())
	}
	response.SetUrl(urlDatas)

	if len(urlDatas) == 0 {
		return &response, status.Error(codes.ResourceExhausted, "no content")
	}
	return &response, nil
}
