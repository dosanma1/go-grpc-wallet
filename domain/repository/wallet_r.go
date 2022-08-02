package repository

import (
	"context"

	"github.com/dosanma1/bluelabs_assessment/domain/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletRepository interface {
	Create(wallet *models.Wallet) error
	Get(userID string) (wallet *models.Wallet, err error)
	Update(userID string, amount int64) (wallet *models.Wallet, err error)
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

func (r *walletRepository) Update(userID string, amount int64) (*models.Wallet, error) {

	tx, err := r.db.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	var wallet models.Wallet
	err = r.db.QueryRow(context.Background(), "SELECT user_id, balance FROM wallet WHERE user_id = $1", userID).Scan(&wallet.UserID, &wallet.Balance)
	if err != nil && err.Error() == pgx.ErrNoRows.Error() {
		return nil, err
	}
	if (wallet.Balance + amount) < 0 {
		return nil, status.Error(codes.PermissionDenied, "insufficient funds")
	}

	err = tx.QueryRow(context.Background(), "UPDATE wallet SET balance = balance + $1 WHERE user_id = $2 RETURNING user_id, balance", amount, userID).Scan(&wallet.UserID, &wallet.Balance)
	if err != nil {
		return nil, err
	}

	return &wallet, err
}
