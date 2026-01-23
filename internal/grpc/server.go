package grpc

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/domurdoc/shortener/api"
	"github.com/domurdoc/shortener/internal/auth"
)

type Server struct {
	port string
	srv  *grpc.Server
	log  *zap.SugaredLogger
}

func NewServer(s *ShortenerServer, log *zap.SugaredLogger, auth *auth.Auth, port int) *Server {
	server := &Server{
		port: fmt.Sprintf(":%d", port),
		srv:  grpc.NewServer(grpc.UnaryInterceptor(CreateAuthIntercepter(auth))),
		log:  log,
	}
	api.RegisterShortenerServiceServer(server.srv, s)
	return server
}

func (s *Server) Start() error {
	s.log.Infow("GRPC server is starting...", "port", s.port)
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		s.log.Warnw("GRPC server failed to listen on port", "port", s.port, "err", err)
		return err
	}
	if err := s.srv.Serve(lis); err != nil {
		s.log.Warnw("GRPC server failed to Serve", "err", err)
		return err
	}
	return nil
}

func (s *Server) Close() error {
	s.srv.GracefulStop()
	s.log.Infow("GRPC server is closed")
	return nil
}
