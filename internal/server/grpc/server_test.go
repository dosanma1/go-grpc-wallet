package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/dosanma1/bluelabs_assessment/config"
	"github.com/dosanma1/bluelabs_assessment/domain/repository"
	"github.com/dosanma1/bluelabs_assessment/domain/service"
	"github.com/dosanma1/bluelabs_assessment/internal/postgresql"
	"github.com/dosanma1/bluelabs_assessment/pkg/pb"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var listener *bufconn.Listener

func init() {
	var logger = logrus.New()
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.0000",
	}
	logger.SetFormatter(formatter)
	logger.SetLevel(logrus.DebugLevel)

	cfg, err := config.ParseConfig()
	if err != nil {
		logger.Fatal(err)
	}

	pgxPool, err := postgresql.NewPgxConn(cfg)
	if err != nil {
		logger.Fatalf("NewPgxConn: %+v", err)
	}

	listener = bufconn.Listen(bufSize)
	walletRepository := repository.NewWalletRepository(pgxPool, logger)
	walletService := service.NewWalletService(walletRepository)

	srv := grpc.NewServer()
	go func() {
		if err := srv.Serve(listener); err != nil {
			logger.Panic(err)
		}
		defer func() {
			srv.GracefulStop()
			pgxPool.Close()
		}()
	}()

	pb.RegisterWalletServiceServer(srv, walletService)
	logger.Println("GRPC Test Server started!")
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return listener.Dial()
}
func TestCreateWallet(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	assert.NoError(t, err)

	serviceClient := pb.NewWalletServiceClient(conn)

	wallet := &pb.Wallet{
		UserId:  uuid.New().String(),
		Balance: 100,
	}

	res, err := serviceClient.CreateWallet(ctx, wallet)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	found, err := serviceClient.GetWallet(ctx, &pb.WalletReq{UserId: wallet.UserId})
	assert.NoError(t, err)
	assert.NotNil(t, found)
}

func TestGetWallet(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	assert.NoError(t, err)

	serviceClient := pb.NewWalletServiceClient(conn)

	wallet := &pb.Wallet{
		UserId:  uuid.New().String(),
		Balance: 100,
	}

	res, err := serviceClient.CreateWallet(ctx, wallet)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	found, err := serviceClient.GetWallet(ctx, &pb.WalletReq{UserId: wallet.UserId})
	assert.NoError(t, err)
	assert.Equal(t, wallet.UserId, found.UserId)
	assert.Equal(t, wallet.Balance, found.Balance)

	found, err = serviceClient.GetWallet(ctx, &pb.WalletReq{UserId: uuid.New().String()})
	assert.Error(t, err)
	if e, ok := status.FromError(err); ok {
		assert.Equal(t, e.Message(), pgx.ErrNoRows.Error())
	} else {
		t.Fail()
	}
	assert.Nil(t, found)
}

func TestUpdateWallet(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	assert.NoError(t, err)

	serviceClient := pb.NewWalletServiceClient(conn)

	wallet := &pb.Wallet{
		UserId:  uuid.New().String(),
		Balance: 100,
	}

	res, err := serviceClient.CreateWallet(ctx, wallet)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	found, err := serviceClient.GetWallet(ctx, &pb.WalletReq{UserId: wallet.UserId})
	assert.NoError(t, err)
	assert.NotNil(t, found)

	var amount float64 = 1000
	found, err = serviceClient.UpdateWallet(ctx, &pb.FundsReq{UserId: wallet.UserId, Amount: amount})
	assert.NoError(t, err)
	assert.Equal(t, wallet.Balance+amount, found.Balance)
}
