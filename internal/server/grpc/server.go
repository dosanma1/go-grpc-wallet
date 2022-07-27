package grpc

import (
	"net"

	"github.com/dosanma1/bluelabs_assessment/config"
	"github.com/dosanma1/bluelabs_assessment/domain/repository"
	"github.com/dosanma1/bluelabs_assessment/domain/service"
	"github.com/dosanma1/bluelabs_assessment/internal/server"
	"github.com/dosanma1/bluelabs_assessment/pkg/pb"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type grpcServer struct {
	log     *logrus.Logger
	cfg     *config.Config
	pgxPool *pgxpool.Pool
}

func NewGrpcServer(log *logrus.Logger, cfg *config.Config, pgxPool *pgxpool.Pool) server.Server {
	return &grpcServer{
		log:     log,
		cfg:     cfg,
		pgxPool: pgxPool,
	}
}

func (s *grpcServer) Run() {
	listener, err := net.Listen(s.cfg.GRPC.Protocol, s.cfg.GRPC.Port)
	if err != nil {
		s.log.Panic(err)
	}

	walletRepository := repository.NewWalletRepository(s.pgxPool, s.log)
	walletService := service.NewWalletService(walletRepository)

	srv := grpc.NewServer()
	pb.RegisterWalletServiceServer(srv, walletService)

	go func() {
		if err := srv.Serve(listener); err != nil {
			s.log.Panic(err)
		}
		defer srv.GracefulStop()
	}()

	s.log.Println("GRPC Server started!")
}
