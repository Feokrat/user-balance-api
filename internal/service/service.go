package service

import (
	"log"

	"github.com/Feokrat/user-balance-api/internal/model"

	"github.com/Feokrat/user-balance-api/internal/repository"
	"github.com/google/uuid"
)

type UserBalance interface {
	GetBalanceByUserId(userId uuid.UUID) (float64, error)
	ChangeUserBalanceByUserId(userId uuid.UUID, changeAmount float64) (bool, error)
	ApplyTransaction(senderId uuid.UUID, receiverId uuid.UUID, amount float64) error
	GetExchangeRate(fromCurrency string, toCurrency string) (float64, error)
}

type TransactionLog interface {
	GetAllUserLogs(userId uuid.UUID, sortField string, pageNum int, pageSize int) ([]model.TransactionLog, error)
	CountUserLogs(userId uuid.UUID) (int, error)
}

type Services struct {
	UserBalance
	TransactionLog
}

func NewServices(repos *repository.Repository, logger *log.Logger) *Services {
	return &Services{
		UserBalance:    NewUserBalanceService(repos.UserBalance, repos.TransactionLog, logger),
		TransactionLog: NewTransactionLogService(repos.TransactionLog, logger),
	}
}
