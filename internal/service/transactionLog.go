package service

import (
	"log"

	"github.com/Feokrat/user-balance-api/internal/model"
	"github.com/Feokrat/user-balance-api/internal/repository"
	"github.com/google/uuid"
)

type TransactionLogService struct {
	transactionLogRepo repository.TransactionLog
	logger             *log.Logger
}

func NewTransactionLogService(transactionLogRepo repository.TransactionLog, logger *log.Logger) *TransactionLogService {
	return &TransactionLogService{transactionLogRepo: transactionLogRepo, logger: logger}
}

func (t TransactionLogService) GetAllUserLogs(userId uuid.UUID, sortField string, pageNum int,
	pageSize int) ([]model.TransactionLog, error) {
	transactionLogs, err := t.transactionLogRepo.GetAllByUserId(userId, sortField, pageNum, pageSize)
	if err != nil {
		t.logger.Printf("could not get all transaction logs of user %v", userId)
		return nil, err
	}

	return transactionLogs, nil
}

func (t TransactionLogService) CountUserLogs(userId uuid.UUID) (int, error) {
	count, err := t.transactionLogRepo.CountByUserId(userId)
	if err != nil {
		t.logger.Printf("could not count transaction logs of user %v, error: %s", userId, err.Error())
		return 0, err
	}

	return count, nil
}