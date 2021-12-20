package internalgrpc

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/config"
	pb "github.com/leksss/hw-test/hw12_13_14_15_calendar/proto/protobuf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcAddr string
	grpc     *grpc.Server
	log      interfaces.Log
	storage  interfaces.Storage
}

func NewServer(log interfaces.Log, config *config.Config, storage interfaces.Storage) *Server {
	return &Server{
		log:      log,
		storage:  storage,
		grpcAddr: config.GRPCAddr.DSN(),
	}
}

func (s *Server) StartGRPC() error {
	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		s.log.Error("failed to listen:", zap.Error(err))
	}
	s.grpc = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(s.log.GetLogger()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(s.log.GetLogger()),
		)),
	)
	pb.RegisterEventServiceServer(s.grpc, NewEventService(s.storage, s.log))

	s.log.Info(fmt.Sprintf("serving gRPC on %s", s.grpcAddr))
	return s.grpc.Serve(lis)
}

func (s *Server) StopGRPC() {
	s.log.Info("stopping gRPC server...")
	s.grpc.GracefulStop()
}
