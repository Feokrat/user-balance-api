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

type Repository struct {
	UserBalance
}

func NewRepositories(db *sqlx.DB, logger *log.Logger) *Repository {
	return &Repository{
		UserBalance: NewUserBalancePostgres(db, logger),
	}
}
