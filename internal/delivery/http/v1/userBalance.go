package v1

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/Feokrat/user-balance-api/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) initUserBalanceRoutes(api *gin.RouterGroup) {
	userBalances := api.Group("/balances")
	{
		userBalances.GET("/:id", h.getUserBalance)
		userBalances.PUT("/", h.changeUserBalance)
		userBalances.POST("/send/", h.sendMoneyFromUserToUser)
		userBalances.GET("/transactionLogs/:id", h.getTransactionLogs)
	}
}

func (h Handler) getTransactionLogs(ctx *gin.Context) {
	userIdStr := ctx.Param("id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})
		return
	}
	pageNum := 1
	pageSize := 1000

	pageNumStr := ctx.Query("pageNum")
	if pageNumStr != "" {
		pageNum, err = strconv.Atoi(pageNumStr)
		if err != nil {
			h.logger.Printf("could not convert pageNum param to int")
			ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
				Message: "could not convert pageNum param to int",
				Errors:  err.Error(),
			})
			return
		}
	}
	pageSizeStr := ctx.Query("pageSize")
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			h.logger.Printf("could not convert pageSize param to int")
			ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
				Message: "could not convert pageSize param to int",
				Errors:  err.Error(),
			})
			return
		}
	}
	sortField := ctx.Query("sortField")
	if sortField == "" {
		sortField = "date"
	}

	logs, err := h.services.GetAllUserLogs(userId, sortField, pageNum-1, pageSize)
	if err != nil {
		h.logger.Printf("could not get all transaction logs of user %v",
			userId)
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, logs)
}

func (h Handler) getUserBalance(ctx *gin.Context) {
	userIdStr := ctx.Param("id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	userBalance, err := h.services.GetBalanceByUserId(userId)
	if err != nil {
		h.logger.Printf("could not get balance of user %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	currencyConvert := ctx.Query("currency")
	if currencyConvert == "" {
		ctx.JSON(http.StatusOK, schemas.UserBalanceResponse{Balance: userBalance})
	} else {
		exchangeRate, err := h.services.GetExchangeRate("", currencyConvert)
		if err != nil {
			h.logger.Printf("could not get exchange rates, error: %s",
				err.Error())
			ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		h.logger.Println(userBalance * exchangeRate)
		ctx.JSON(http.StatusOK, schemas.UserBalanceResponse{
			Balance: math.Ceil(userBalance*exchangeRate*100) / 100,
		})
		return
	}
}

func (h Handler) changeUserBalance(ctx *gin.Context) {
	var requestModel schemas.ChangeBalanceRequest

	if err := ctx.BindJSON(&requestModel); err != nil {
		h.logger.Printf("request body in wrong format, error: %s",
			err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong request format",
			Errors:  err.Error(),
		})
		return
	}

	created, err := h.services.ChangeUserBalanceByUserId(requestModel.UserId, requestModel.ChangeAmount)
	if err != nil {
		h.logger.Printf("could not change balance of user %v, error: %s",
			requestModel.UserId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if created {
		h.logger.Printf("created a new account balance")
		ctx.JSON(http.StatusCreated, "")
	}
}

func (h Handler) sendMoneyFromUserToUser(ctx *gin.Context) {
	var requestModel schemas.TransactionRequest

	if err := ctx.BindJSON(&requestModel); err != nil {
		h.logger.Printf("request body in wrong format, error: %s",
			err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong request model",
			Errors:  err.Error(),
		})
		return
	}

	if requestModel.Amount < 0 {
		err := schemas.ErrorAmountToSendNegative{
			Message: fmt.Sprintf("amount of sending money is negative: %v < 0",
				requestModel.Amount),
		}
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "amount of sending money can not be negative",
			Errors:  err.Error(),
		})
		return
	}

	err := h.services.ApplyTransaction(requestModel.SenderId, requestModel.ReceiverId, requestModel.Amount)
	if err != nil {
		h.logger.Printf("could not apply transaction from user %v to user %v, error: %s",
			requestModel.SenderId, requestModel.ReceiverId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})
		return
	}
}
