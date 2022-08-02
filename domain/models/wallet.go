package models

import (
	"log"

	"github.com/dosanma1/bluelabs_assessment/pkg/pb"
	"github.com/google/uuid"
)

type Wallet struct {
	UserID  uuid.UUID `json:"user_id"`
	Balance int64     `json:"balance"`
}

func (w *Wallet) ToProtoBuffer() *pb.Wallet {
	return &pb.Wallet{
		UserId:  w.UserID.String(),
		Balance: w.Balance,
	}
}

func (w *Wallet) FromProtoBuffer(wallet *pb.Wallet) {
	_uuid, err := uuid.Parse(wallet.GetUserId())
	if err != nil {
		log.Fatalln(_uuid)
	}
	w.UserID = _uuid
	w.Balance = wallet.GetBalance()
}
