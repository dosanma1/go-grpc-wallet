package repository

import (
	"context"
	"testing"

	"github.com/dosanma1/go-grpc-wallet/config"
	"github.com/dosanma1/go-grpc-wallet/domain/models"
	"github.com/dosanma1/go-grpc-wallet/internal/postgresql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	var logger = logrus.New()

	cfg, err := config.ParseConfig()
	if err != nil {
		logger.Fatal(err)
	}

	pgxPool, err := postgresql.NewPgxConn(cfg)
	if err != nil {
		logger.Fatalf("NewPgxConn: %+v", err)
	}
	defer pgxPool.Close()

	_, err = pgxPool.Exec(context.Background(), "DELETE FROM wallet")
	if err != nil {
		logger.Panic(err)
	}
}

func TestWalletRepositoryCreate(t *testing.T) {
	var logger = logrus.New()

	cfg, err := config.ParseConfig()
	if err != nil {
		logger.Fatal(err)
	}

	pgxPool, err := postgresql.NewPgxConn(cfg)
	if err != nil {
		logger.Fatalf("NewPgxConn: %+v", err)
	}
	assert.NoError(t, err)
	defer pgxPool.Close()

	wallet := &models.Wallet{
		UserID:  uuid.New(),
		Balance: 0,
	}

	r := NewWalletRepository(pgxPool, logger)
	r.Create(wallet)
	assert.NoError(t, err)

	found, err := r.Get(wallet.UserID.String())
	assert.NoError(t, err)
	assert.NotNil(t, found)
}

func TestWalletRepositoryGet(t *testing.T) {
	var logger = logrus.New()

	cfg, err := config.ParseConfig()
	if err != nil {
		logger.Fatal(err)
	}

	pgxPool, err := postgresql.NewPgxConn(cfg)
	if err != nil {
		logger.Fatalf("NewPgxConn: %+v", err)
	}
	assert.NoError(t, err)
	defer pgxPool.Close()

	wallet := &models.Wallet{
		UserID:  uuid.New(),
		Balance: 0,
	}

	r := NewWalletRepository(pgxPool, logger)
	r.Create(wallet)
	assert.NoError(t, err)

	found, err := r.Get(wallet.UserID.String())
	assert.NoError(t, err)
	assert.Equal(t, wallet.UserID, found.UserID)
	assert.Equal(t, wallet.Balance, found.Balance)

	found, err = r.Get(uuid.New().String())
	assert.Error(t, err)
	assert.EqualError(t, pgx.ErrNoRows, err.Error())
	assert.Nil(t, found)
}

func TestWalletRepositoryUpdate(t *testing.T) {
	var logger = logrus.New()

	cfg, err := config.ParseConfig()
	if err != nil {
		logger.Fatal(err)
	}

	pgxPool, err := postgresql.NewPgxConn(cfg)
	if err != nil {
		logger.Fatalf("NewPgxConn: %+v", err)
	}
	assert.NoError(t, err)
	defer pgxPool.Close()

	wallet := &models.Wallet{
		UserID:  uuid.New(),
		Balance: 0,
	}

	r := NewWalletRepository(pgxPool, logger)
	r.Create(wallet)
	assert.NoError(t, err)

	found, err := r.Get(wallet.UserID.String())
	assert.NoError(t, err)
	assert.NotNil(t, found)

	var amount int64 = 200
	updated, err := r.Update(wallet.UserID.String(), amount)
	assert.NoError(t, err)
	assert.NotNil(t, updated)

	found, err = r.Get(wallet.UserID.String())
	assert.NoError(t, err)
	assert.Equal(t, wallet.Balance+amount, found.Balance)
}
