package repository

import (
	"database/sql"
	"github.com/Feokrat/user-balance-api/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"os"
	"testing"
)

func TestUserBalancePostgres_Create(t *testing.T) {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("error while trying to mock db, error: %s", err.Error())
	}
	defer db.Close()

	r := NewUserBalancePostgres(db, logger)

	type args struct {
		userBalance model.UserBalance
	}

	type mockBehavior func(args args)

	testUserId := uuid.New()

	tests := []struct {
		name        string
		mock        mockBehavior
		input       args
		expectedErr bool
	}{
		{
			name: "Ok",
			mock: func(args args) {
				userBalance := args.userBalance
				rows := sqlxmock.NewRows([]string{"user_id"}).AddRow(testUserId)
				mock.ExpectQuery("INSERT INTO user_balance").
					WithArgs(userBalance.UserId, userBalance.Balance).
					WillReturnRows(rows)
			},
			input: args{userBalance: model.UserBalance{
				UserId:  testUserId,
				Balance: 100,
			}},
			expectedErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.input)

			err := r.Create(test.input.userBalance)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserBalancePostgres_GetByUserId(t *testing.T) {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)

	db, mock, err := sqlxmock.Newx(sqlxmock.QueryMatcherOption(sqlxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error while trying to mock db, error: %s", err.Error())
	}
	defer db.Close()

	r := NewUserBalancePostgres(db, logger)

	type args struct {
		userId uuid.UUID
	}

	testUserId := uuid.New()

	type mockBehavior func(args args)

	tests := []struct {
		name        string
		mock        mockBehavior
		input       args
		expectedOut model.UserBalance
		expectedErr bool
		err error
	}{
		{
			name: "Ok",
			mock: func(args args) {
				rows := sqlxmock.NewRows([]string{"user_id", "balance"}).
					AddRow(testUserId, 20)

				mock.ExpectQuery("SELECT ub.user_id, ub.balance FROM user_balance AS ub WHERE ub.user_id = $1").
					WithArgs(args.userId).WillReturnRows(rows)
			},
			input: args{userId: testUserId},
			expectedOut: model.UserBalance{
				UserId:  testUserId,
				Balance: 20,
			},
			expectedErr: false,
			err: nil,
		},
		{
			name: "Not found",
			mock: func(args args) {
				rows := sqlxmock.NewRows([]string{"user_id", "balance"})

				mock.ExpectQuery("SELECT ub.user_id, ub.balance FROM user_balance AS ub WHERE ub.user_id = $1").
					WithArgs(args.userId).WillReturnRows(rows)
			},
			input: args{userId: testUserId},
			expectedOut: model.UserBalance{},
			expectedErr: true,
			err: sql.ErrNoRows,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.input)

			got, err := r.GetByUserId(test.input.userId)
			if test.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, test.err)
			} else {
				assert.Equal(t, got.UserId, test.expectedOut.UserId)
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserBalancePostgres_UpdateByUserId(t *testing.T) {

}