package repository

import (
	"context"

	"github.com/dosanma1/bluelabs_assessment/domain/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type WalletRepository interface {
	Create(wallet *models.Wallet) error
	Get(userID string) (wallet *models.Wallet, err error)
	Update(userID string, amount int64) error
}

type walletRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewWalletRepository(db *pgxpool.Pool, log *logrus.Logger) *walletRepository {
	return &walletRepository{
		db:  db,
		log: log,
	}
}

func (r *walletRepository) Create(wallet *models.Wallet) error {
	_, err := r.db.Exec(context.Background(), "INSERT INTO wallet (user_id, balance) VALUES ($1, $2)", wallet.UserID, wallet.Balance)
	return err
}

func (r *walletRepository) Get(userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.QueryRow(context.Background(), "SELECT user_id, balance FROM wallet WHERE user_id = $1", userID).Scan(&wallet.UserID, &wallet.Balance)
	if err != nil && err.Error() == pgx.ErrNoRows.Error() {
		return nil, err
	}

	return &wallet, err
}

func (r *walletRepository) Update(userID string, amount int64) error {
	_, err := r.db.Exec(context.Background(), "UPDATE wallet SET balance = balance + $1 WHERE user_id = $2 RETURNING balance", amount, userID)
	if err != nil {
		return err
	}

	return err
}
