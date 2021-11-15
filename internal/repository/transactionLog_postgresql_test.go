package repository

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"

	"github.com/Feokrat/user-balance-api/internal/model"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTransactionLogPostgres_Create(t *testing.T) {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("error while trying to mock db, error: %s", err.Error())
	}
	defer db.Close()

	r := NewTransactionLogPostgres(db, logger)

	type args struct {
		transactionLog model.TransactionLog
	}

	type mockBehavior func(args args)

	testUserId := uuid.New()

	tests := []struct {
		name        string
		mock        mockBehavior
		input       args
		expectedOut int32
		expectedErr bool
	}{
		{
			name: "Ok",
			input: args{transactionLog: model.TransactionLog{
				UserId:     testUserId,
				Date:       time.Now(),
				Amount:     100,
				Commentary: "Test 100",
			}},
			mock: func(args args) {
				transactionLog := args.transactionLog
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO transaction_log").
					WithArgs(transactionLog.UserId, transactionLog.Date, transactionLog.Amount, transactionLog.Commentary).
					WillReturnRows(rows)
			},
			expectedOut: 1,
			expectedErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.input)

			got, err := r.Create(test.input.transactionLog)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, got)
			}
		})
	}
}

func TestTransactionLogPostgres_GetAllByUserId(t *testing.T) {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)

	db, mock, err := sqlxmock.Newx(sqlxmock.QueryMatcherOption(sqlxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error while trying to mock db, error: %s", err.Error())
	}
	defer db.Close()

	r := NewTransactionLogPostgres(db, logger)

	type args struct {
		userId uuid.UUID
		pageNum int
		pageSize int
	}

	type mockBehavior func(args args)

	userId := uuid.New()
	time := time.Now()

	tests := []struct {
		name        string
		mock        mockBehavior
		input       args
		expectedOut []model.TransactionLog
		expectedErr bool
	}{
		{
			name:  "Ok",
			input: args{
				userId:   userId,
				pageNum:  1,
				pageSize: 100,
			},
			mock: func(args args) {
				rows := sqlxmock.NewRows([]string{"id", "user_id", "date", "amount", "commentary"}).
					AddRow(1, userId, time, 100, "TEST1").
					AddRow(2, userId, time, 200, "TEST2")

				mock.ExpectQuery("SELECT tl.id, tl.user_id, tl.date, tl.amount, tl.commentary FROM transaction_log AS tl WHERE tl.user_id = $1 ORDER BY date LIMIT $2 OFFSET $3").
					WithArgs(args.userId, args.pageSize, args.pageNum*args.pageSize).WillReturnRows(rows)
			},
			expectedOut: []model.TransactionLog{
				{
					Id:         1,
					UserId:     userId,
					Date:       time,
					Amount:     100,
					Commentary: "TEST1",
				},
				{
					Id:         2,
					UserId:     userId,
					Date:       time,
					Amount:     200,
					Commentary: "TEST2",
				},
			},
			expectedErr: false,
		},
		{
			name:  "Empty List",
			input: args{
				userId:   userId,
				pageNum:  1,
				pageSize: 100,
			},
			mock: func(args args) {
				rows := sqlxmock.NewRows([]string{"id", "user_id", "date", "amount", "commentary"})

				mock.ExpectQuery("SELECT tl.id, tl.user_id, tl.date, tl.amount, tl.commentary FROM transaction_log AS tl WHERE tl.user_id = $1 ORDER BY date LIMIT $2 OFFSET $3").
					WithArgs(args.userId, args.pageSize, args.pageNum*args.pageSize).WillReturnRows(rows)
			},
			expectedOut: nil,
			expectedErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.input)

			got, err := r.GetAllByUserId(test.input.userId, "date", 1, 100)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOut, got)
			}
		})
	}
}