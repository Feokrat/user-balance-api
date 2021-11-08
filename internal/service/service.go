package service

import (
	"log"

	"github.com/Feokrat/user-balance-api/internal/repository"
	"github.com/google/uuid"
)

type UserBalance interface {
	GetBalanceByUserId(userId uuid.UUID) (float64, error)
	ChangeUserBalanceByUserId(userId uuid.UUID, changeAmount float64) (bool, error)
	ApplyTransaction(senderId uuid.UUID, receiverId uuid.UUID, amount float64) error
	GetExchangeRate(fromCurrency string, toCurrency string) (float64, error)
}

type Service struct {
	UserBalance
}

func NewServices(repos *repository.Repository, logger *log.Logger) *Service {
	return &Service{
		UserBalance: NewUserBalanceService(repos.UserBalance, logger),
	}
}
