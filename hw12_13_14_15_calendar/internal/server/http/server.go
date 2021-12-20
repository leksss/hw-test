package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/config"
	pb "github.com/leksss/hw-test/hw12_13_14_15_calendar/proto/protobuf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcAddr string
	http     *http.Server
	log      interfaces.Log
}

func NewServer(log interfaces.Log, config *config.Config, storage interfaces.Storage) *Server {
	return &Server{
		log:      log,
		grpcAddr: config.GRPCAddr.DSN(),
		http: &http.Server{
			Addr: config.HTTPAddr.DSN(),
		},
	}
}

func (s *Server) StartHTTPProxy() error {
	conn, err := grpc.DialContext(
		context.Background(),
		s.grpcAddr,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		s.log.Error("failed to dial server:", zap.Error(err))
	}

	gwMux := runtime.NewServeMux()
	err = pb.RegisterEventServiceHandler(context.Background(), gwMux, conn)
	if err != nil {
		s.log.Error("failed to register gateway:", zap.Error(err))
	}

	s.http.Handler = loggingMiddleware(gwMux, s.log)
	s.log.Info(fmt.Sprintf("serving gRPC-Gateway on %s", s.http.Addr))
	return s.http.ListenAndServe()
}

func (s *Server) StopHTTPProxy(ctx context.Context) error {
	s.log.Info("stopping HTTP proxy server...")
	return s.http.Shutdown(ctx)
}
