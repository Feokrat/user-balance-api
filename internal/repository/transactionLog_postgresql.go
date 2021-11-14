package repository

import (
	"fmt"
	"log"

	"github.com/Feokrat/user-balance-api/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func (t TransactionLogPostgres) GetAllByUserId(userId uuid.UUID, sortField string, pageNum int, pageSize int) ([]model.TransactionLog, error) {
	query := fmt.Sprintf("SELECT tl.id, tl.user_id, tl.date, tl.amount, tl.commentary FROM transaction_log AS tl "+
		"WHERE tl.user_id = $1 ORDER BY %s LIMIT $2 OFFSET $3", sortField)

	var transactionLogs []model.TransactionLog

	err := t.db.Select(&transactionLogs, query, userId, pageSize, pageNum*pageSize)
	if err != nil {
		t.logger.Printf("error in db while trying to get transaction log of user %v, error: %s",
			userId, err)
		return nil, err
	}

	return transactionLogs, nil
}

func (t TransactionLogPostgres) Create(transactionLog model.TransactionLog) (int32, error) {
	query := "INSERT INTO transaction_log AS tl (user_id, date, amount, commentary) VALUES ($1, $2, $3, $4) RETURNING id"

	var id int32

	row := t.db.QueryRow(query, transactionLog.UserId, transactionLog.Date, transactionLog.Amount, transactionLog.Commentary)

	if err := row.Scan(&id); err != nil {
		t.logger.Printf("error in db while trying to create transaction log info for user %v, error: %s",
			transactionLog.UserId, err.Error())
		return 0, err
	}

	return id, nil
}
