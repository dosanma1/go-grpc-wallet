package service

import (
	"context"

	"github.com/dosanma1/bluelabs_assessment/domain/models"
	"github.com/dosanma1/bluelabs_assessment/domain/repository"
	"github.com/dosanma1/bluelabs_assessment/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type walletService struct {
	pb.UnimplementedWalletServiceServer

	walletRepository repository.WalletRepository
}

func NewWalletService(repository repository.WalletRepository) pb.WalletServiceServer {
	return &walletService{
		walletRepository: repository,
	}
}

func (s *walletService) CreateWallet(ctx context.Context, req *pb.Wallet) (*pb.Wallet, error) {
	var wallet models.Wallet
	wallet.FromProtoBuffer(req)
	err := s.walletRepository.Create(&wallet)
	if err != nil {
		return nil, status.FromContextError(err).Err()
	}

	return wallet.ToProtoBuffer(), nil
}
func (s *walletService) GetWallet(ctx context.Context, req *pb.WalletReq) (*pb.Wallet, error) {
	wallet, err := s.walletRepository.Get(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return wallet.ToProtoBuffer(), err
}

func (s *walletService) UpdateWallet(ctx context.Context, req *pb.FundsReq) (*pb.Wallet, error) {

	wallet, err := s.walletRepository.Get(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if (wallet.Balance + req.Amount) < 0 {
		return nil, status.Error(codes.PermissionDenied, "insufficient funds")
	}

	err = s.walletRepository.Update(req.UserId, req.Amount)
	if err != nil {
		return nil, err
	}

	wallet, err = s.walletRepository.Get(req.UserId)

	return wallet.ToProtoBuffer(), err
}
