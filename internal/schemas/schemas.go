package schemas

import "github.com/google/uuid"

type ErrorResponse struct {
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

type ErrorAmountToSendNegative struct {
	Message string `json:"message"`
}

func (e ErrorAmountToSendNegative) Error() string {
	return e.Message
}

type ErrorNotEnoughFunds struct {
	Message string `json:"message"`
}

func (e ErrorNotEnoughFunds) Error() string {
	return e.Message
}

type ErrorUserBalanceNotFound struct {
	Message string `json:"message"`
}

func (e ErrorUserBalanceNotFound) Error() string {
	return e.Message
}

type ValidationErrorResponse struct {
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type UserBalanceResponse struct {
	Balance float64 `json:"balance"`
}

type ChangeBalanceRequest struct {
	UserId       uuid.UUID `json:"userId"`
	ChangeAmount float64   `json:"changeAmount"`
}

type TransactionRequest struct {
	SenderId   uuid.UUID `json:"senderId"`
	ReceiverId uuid.UUID `json:"receiverId"`
	Amount     float64   `json:"amount"`
}
