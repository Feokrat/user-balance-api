package repository

import (
	"log"

	"github.com/Feokrat/user-balance-api/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserBalance interface {
	GetByUserId(userId uuid.UUID) (model.UserBalance, error)
	UpdateByUserId(userId uuid.UUID, changeAmount float64) error
	CheckIfExistsByUserId(userId uuid.UUID) (bool, error)
	Create(userBalance model.UserBalance) error
}

type TransactionLog interface {
	GetAllByUserId(userId uuid.UUID, sortField string, pageNum int, pageSize int) ([]model.TransactionLog, error)
	Create(transactionLog model.TransactionLog) (int32, error)
}

type Repository struct {
	UserBalance
	TransactionLog
}

func NewRepositories(db *sqlx.DB, logger *log.Logger) *Repository {
	return &Repository{
		UserBalance:    NewUserBalancePostgres(db, logger),
		TransactionLog: NewTransactionLogPostgres(db, logger),
	}
}
