package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Feokrat/user-balance-api/internal/schemas"

	"github.com/Feokrat/user-balance-api/internal/model"

	"github.com/google/uuid"

	"github.com/Feokrat/user-balance-api/internal/repository"
)

const (
	ECHANGE_RATE_URL      = "https://free.currconv.com/api/v7/convert?"
	EXCHANGE_RATE_API_KEY = "620c10be1a7ef027dcd9"
	BASE_CURRENCY         = "RUB"
)

type UserBalanceService struct {
	userBalanceRepo    repository.UserBalance
	transactionLogRepo repository.TransactionLog
	logger             *log.Logger
}

func NewUserBalanceService(userBalanceRepo repository.UserBalance, transactionLogRepo repository.TransactionLog,
	logger *log.Logger) *UserBalanceService {
	return &UserBalanceService{userBalanceRepo: userBalanceRepo, transactionLogRepo: transactionLogRepo, logger: logger}
}

func (s UserBalanceService) GetBalanceByUserId(userId uuid.UUID) (float64, error) {
	ubExists, err := s.userBalanceRepo.CheckIfExistsByUserId(userId)
	if err != nil {
		s.logger.Printf("could not check if user with id %v exists, error: %s",
			userId, err.Error())
		return 0, err
	}

	if !ubExists {
		s.logger.Printf("user with id %v does not exist",
			userId)
		return 0, nil
	}

	ub, err := s.userBalanceRepo.GetByUserId(userId)
	if err != nil {
		s.logger.Printf("could not get user balance info of user with id %v, error: %s",
			userId, err.Error())
		return 0, err
	}

	return ub.Balance, nil
}

func (s UserBalanceService) ChangeUserBalanceByUserId(userId uuid.UUID, changeAmount float64) (bool, error) {
	ubExists, err := s.userBalanceRepo.CheckIfExistsByUserId(userId)
	if err != nil {
		s.logger.Printf("could not check if user %v exists, error: %s",
			userId, err.Error())
		return false, err
	}

	if changeAmount > 0 {
		s.logger.Printf("trying to add balance to user %v",
			userId)

		if !ubExists {
			s.logger.Printf("user %v does not exist, trying to create him with balance %v",
				userId, changeAmount)

			err = s.userBalanceRepo.Create(model.UserBalance{
				UserId:  userId,
				Balance: changeAmount,
			})
			if err != nil {
				s.logger.Printf("could not create user balance of user %v and balance %v",
					userId, changeAmount)
				return false, err
			}

			err = s.logBalanceInfo(userId, changeAmount, fmt.Sprintf("Added %v rubles",
				changeAmount))
			if err != nil {
				s.logger.Printf("could not log info about user %v, error: %s",
					userId, err.Error())
				return false, err
			}

			return true, nil
		} else {
			err = s.logBalanceInfo(userId, changeAmount, fmt.Sprintf("Added %v rubles",
				changeAmount))
			if err != nil {
				s.logger.Printf("could not log info about user %v, error: %s",
					userId, err.Error())
				return false, err
			}

			return false, s.addBalance(userId, changeAmount)
		}
	} else {
		s.logger.Printf("trying to sub balance of user %v",
			userId)

		if !ubExists {
			s.logger.Printf("user %v does not exist to sub his balance",
				userId)
			return false, schemas.ErrorUserBalanceNotFound{
				Message: fmt.Sprintf("user balance of user with id %v not found",
					userId),
			}
		} else {
			err = s.logBalanceInfo(userId, changeAmount, fmt.Sprintf("Substracted %v rubles",
				math.Abs(changeAmount)))
			if err != nil {
				s.logger.Printf("could not log info about user %v, error: %s",
					userId, err.Error())
				return false, err
			}

			return false, s.subBalance(userId, changeAmount)
		}
	}
}

func (s UserBalanceService) ApplyTransaction(senderId uuid.UUID, receiverId uuid.UUID, amount float64) error {
	senderAbExists, err := s.userBalanceRepo.CheckIfExistsByUserId(senderId)
	if err != nil {
		s.logger.Printf("could not check if sender %v exists, error: %s",
			senderId, err.Error())
		return err
	}

	receiverAbExists, err := s.userBalanceRepo.CheckIfExistsByUserId(receiverId)
	if err != nil {
		s.logger.Printf("could not check if receiver %v exists, error: %s",
			senderId, err.Error())
		return err
	}

	if !senderAbExists {
		s.logger.Printf("sender %v does not exist to sub his balance",
			senderId)
		return schemas.ErrorUserBalanceNotFound{
			Message: fmt.Sprintf("user balance of sender %v not found",
				senderId),
		}
	}

	if !receiverAbExists {
		s.logger.Printf("receiver %v does not exist to add to his balance",
			senderId)
		return schemas.ErrorUserBalanceNotFound{
			Message: fmt.Sprintf("user balance of receiver %v not found",
				senderId),
		}
	}

	err = s.subBalance(senderId, -amount)
	if err != nil {
		s.logger.Printf("could not receive money from user %v for transaction to user %v balance, error: %s",
			senderId, receiverId, err.Error())
		return err
	}

	err = s.logBalanceInfo(senderId, amount, fmt.Sprintf("Sended %v rubles to user %v",
		amount, receiverId))
	if err != nil {
		s.logger.Printf("could not log info about user %v, error: %s",
			senderId, err.Error())
		return err
	}

	err = s.addBalance(receiverId, amount)
	if err != nil {
		s.logger.Printf("could not send money to user %v, trying to return money to user %v, error: %v",
			receiverId, senderId, err.Error())
		err = s.addBalance(senderId, amount)
		if err != nil {
			s.logger.Printf("could not return money to user %v, error: %s",
				senderId, err.Error())
		}
		return err
	}

	err = s.logBalanceInfo(receiverId, amount, fmt.Sprintf("Received %v rubles from user %v",
		amount, senderId))
	if err != nil {
		s.logger.Printf("could not log info about user %v, error: %s",
			receiverId, err.Error())
		return err
	}

	return nil
}

func (s UserBalanceService) GetExchangeRate(fromCurrency string, toCurrency string) (float64, error) {
	if fromCurrency == "" {
		fromCurrency = BASE_CURRENCY
	}

	currencies := fmt.Sprintf("%s_%s", fromCurrency, toCurrency)
	exchangerURL := fmt.Sprintf("%sq=%s&compact=ultra&apiKey=%s",
		ECHANGE_RATE_URL, currencies, EXCHANGE_RATE_API_KEY)
	resp, err := http.Get(exchangerURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		s.logger.Printf("could not get exchange rates, status code: %v, error: %s",
			resp.StatusCode, err.Error())
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return 0, err
	}

	respBody := map[string]float64{}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s",
			err.Error())
		return 0, err
	}

	return respBody[currencies], nil
}

func (s UserBalanceService) addBalance(userId uuid.UUID, changeAmount float64) error {
	err := s.userBalanceRepo.UpdateByUserId(userId, changeAmount)
	if err != nil {
		s.logger.Printf("could not add balance to user %v, error: %s",
			userId, err.Error())
		return err
	}

	return nil
}

func (s UserBalanceService) subBalance(userId uuid.UUID, changeAmount float64) error {
	ub, err := s.userBalanceRepo.GetByUserId(userId)
	if err != nil {
		return err
	}
	if math.Abs(changeAmount) > ub.Balance {
		s.logger.Printf("Not enough funds in user %v balance", userId)
		return schemas.ErrorNotEnoughFunds{
			Message: fmt.Sprintf("User %v has less money than %v",
				userId, math.Abs(changeAmount)),
		}
	}
	err = s.userBalanceRepo.UpdateByUserId(userId, changeAmount)
	if err != nil {
		s.logger.Printf("could not sub balance of a user %v, error: %s",
			userId, err.Error())
		return err
	}

	return nil
}

func (s UserBalanceService) logBalanceInfo(userId uuid.UUID, amount float64, commentary string) error {
	_, err := s.transactionLogRepo.Create(model.TransactionLog{
		UserId:     userId,
		Date:       time.Now(),
		Amount:     math.Abs(amount),
		Commentary: commentary,
	})

	return err
}
