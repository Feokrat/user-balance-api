package repository

import (
	"database/sql"
	"log"

	"github.com/Feokrat/user-balance-api/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserBalancePostgres struct {
	db     *sqlx.DB
	logger *log.Logger
}

func NewUserBalancePostgres(db *sqlx.DB, logger *log.Logger) *UserBalancePostgres {
	return &UserBalancePostgres{
		db:     db,
		logger: logger}
}

func (r UserBalancePostgres) GetByUserId(userId uuid.UUID) (model.UserBalance, error) {
	query := "SELECT ub.user_id, ub.balance FROM user_balance AS ub WHERE ub.user_id = $1"

	var userBalance model.UserBalance

	err := r.db.Get(&userBalance, query, userId)

	if err != nil {
		r.logger.Printf("error in db while trying to get user balance of user %v, error: %s",
			userId, userId)
		return model.UserBalance{}, err
	}

	return userBalance, nil
}

func (r UserBalancePostgres) UpdateByUserId(userId uuid.UUID, changeAmount float64) error {
	query := "UPDATE user_balance ub SET balance = balance + $1 WHERE user_id = $2"
	_, err := r.db.Exec(query, changeAmount, userId)
	return err
}

func (r UserBalancePostgres) CheckIfExistsByUserId(userId uuid.UUID) (bool, error) {
	query := "SELECT ub.user_id, ub.balance FROM user_balance AS ub WHERE ub.user_id = $1"

	var UserBalance model.UserBalance

	err := r.db.Get(&UserBalance, query, userId)
	if err == sql.ErrNoRows {
		r.logger.Printf("could not find user balance of user %v in db",
			userId)
		return false, nil
	} else if err != nil {
		r.logger.Printf("error in db while trying to check if user %v exists, error: %s",
			userId, err.Error())
		return false, err
	}

	return true, nil
}

func (r UserBalancePostgres) Create(UserBalance model.UserBalance) error {
	query := "INSERT INTO user_balance AS ub (user_id, balance) VALUES ($1, $2) RETURNING user_id"

	var userId uuid.UUID

	row := r.db.QueryRow(query, UserBalance.UserId, UserBalance.Balance)

	if err := row.Scan(&userId); err != nil {
		r.logger.Printf("error in db while trying to create user balance info of user %v, error: %s",
			userId, err.Error())
		return err
	}

	return nil
}
