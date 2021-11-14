package model

import (
	"github.com/google/uuid"
	"time"
)

type TransactionLog struct {
	Id int32 `json:"id" db:"id"`
	UserId uuid.UUID `json:"userId" db:"user_id"`
	Date time.Time `json:"date" db:"date"`
	Amount float64 `json:"amount" db:"amount"`
	Commentary string `json:"commentary" db:"commentary"`
}
