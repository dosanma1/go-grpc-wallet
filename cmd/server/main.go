package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dosanma1/go-grpc-wallet/config"
	"github.com/dosanma1/go-grpc-wallet/internal/postgresql"
	"github.com/dosanma1/go-grpc-wallet/internal/server/grpc"
	"github.com/dosanma1/go-grpc-wallet/internal/server/http"
	"github.com/sirupsen/logrus"
)

const SERVICE_NAME = "wallet"

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

	pgxPool, err := postgresql.NewPgxConn(cfg)
	if err != nil {
		logger.Fatalf("NewPgxConn: %+v", err)
	}
	defer pgxPool.Close()

	grpcServer := grpc.NewGrpcServer(logger, cfg, pgxPool)
	grpcServer.Run()

	httpServer := http.NewHttpServer(logger, cfg, pgxPool)
	httpServer.Run()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Press ctrl+c to exit")

	go func() {
		_ = <-sigs
		done <- true
	}()

	<-done
	fmt.Printf("%s service exit....\n", SERVICE_NAME)

}
