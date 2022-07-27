package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/dosanma1/bluelabs_assessment/config"
	"github.com/dosanma1/bluelabs_assessment/pkg/pb"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
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

	addr := fmt.Sprintf("%s%s", cfg.GRPC.Host, cfg.GRPC.Port)
	logger.Println(addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic("cannot connect with server " + err.Error())
	}
	defer conn.Close()

	serviceClient := pb.NewWalletServiceClient(conn)

	ctx := context.Background()
	wallet := &pb.Wallet{
		UserId:  uuid.New().String(),
		Balance: 100,
	}

	// CREATE
	res, err := serviceClient.CreateWallet(ctx, wallet)
	if err != nil {
		e, _ := status.FromError(err)
		logger.Printf("wallet not created " + e.Message())
		return
	}
	logger.Printf("Wallet %s created with balance %v", res.UserId, res.Balance)

	// GET
	res, err = serviceClient.GetWallet(ctx, &pb.WalletReq{
		UserId: wallet.UserId,
	})
	if err != nil {
		e, _ := status.FromError(err)
		logger.Printf("wallet not retrieved " + e.Message())
		return
	}
	logger.Printf("Wallet %s balance: %v", res.UserId, res.Balance)

	wg := sync.WaitGroup{}
	// UPDATE
	wg.Add(4)

	go func() {
		res, err = serviceClient.UpdateWallet(ctx, &pb.FundsReq{
			UserId: wallet.UserId,
			Amount: 120,
		})
		if err != nil {
			e, _ := status.FromError(err)
			logger.Printf("wallet not updated: %s", e.Message())
			wg.Done()
			return
		}
		logger.Printf("Updated wallet %s, current balance: %v", res.UserId, res.Balance)
		wg.Done()
	}()
	go func() {
		time.Sleep(500 * time.Millisecond)
		res, err = serviceClient.UpdateWallet(ctx, &pb.FundsReq{
			UserId: wallet.UserId,
			Amount: -200,
		})
		if err != nil {
			e, _ := status.FromError(err)
			logger.Printf("wallet not updated: %s", e.Message())
			wg.Done()
			return
		}
		logger.Printf("Updated wallet %s, current balance: %v", res.UserId, res.Balance)
		wg.Done()
	}()
	go func() {
		time.Sleep(1 * time.Second)
		res, err = serviceClient.UpdateWallet(ctx, &pb.FundsReq{
			UserId: wallet.UserId,
			Amount: -1000,
		})
		if err != nil {
			e, _ := status.FromError(err)
			logger.Printf("wallet not updated: %s", e.Message())
			wg.Done()
			return
		}
		logger.Printf("Updated wallet %s, current balance: %v", res.UserId, res.Balance)
		wg.Done()
	}()
	go func() {
		res, err = serviceClient.UpdateWallet(ctx, &pb.FundsReq{
			UserId: wallet.UserId,
			Amount: 120,
		})
		if err != nil {
			e, _ := status.FromError(err)
			logger.Printf("wallet not updated: %s", e.Message())
			wg.Done()
			return
		}
		logger.Printf("Updated wallet %s, current balance: %v", res.UserId, res.Balance)
		wg.Done()
	}()

	wg.Wait()
}

func generateID() string {
	rand.Seed(time.Now().Unix())
	return "ID: " + strconv.Itoa(rand.Int())
}
