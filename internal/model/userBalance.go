package model

import "github.com/google/uuid"

type UserBalance struct {
	UserId  uuid.UUID `json:"userId" db:"user_id"`
	Balance float64   `json:"balance" db:"balance"`
}
