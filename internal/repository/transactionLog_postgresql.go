package repository

import (
	"github.com/Feokrat/user-balance-api/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
)

type TransactionLogPostgres struct {
	db     *sqlx.DB
	logger *log.Logger
}

func NewTransactionLogPostgres(db *sqlx.DB, logger *log.Logger) *TransactionLogPostgres {
	return &TransactionLogPostgres{
		db:     db,
		logger: logger}
}

func (t TransactionLogPostgres) GetAllUserLogs(userId uuid.UUID, pageNum int, pageSize int) []model.TransactionLog {
	panic("implement me")
}

func (t TransactionLogPostgres) CreateUserLog(transactionLog model.TransactionLog) (int32, error) {
	panic("implement me")
}