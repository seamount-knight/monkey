package control_test

import (
	"github.com/golang/mock/gomock"
	"monkey/mock"
	"testing"
	"monkey/control"
	"monkey/infra/log"
	"monkey/domain"
	er "errors"
	"monkey/infra/errors"
	"github.com/stretchr/testify/assert"
	"github.com/satori/go.uuid"
)

func getStore(ctrl *gomock.Controller) *mock.MockMonkeyStore {
	dbMock := mock.NewMockMonkeyStore(ctrl)
	return dbMock
}

func TestListMonkeys(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := getStore(ctrl)
	controller := control.NewMonkeyController(store, log.NewLogger("control_test"))

	testTable := []struct{
		TestName string
		Prepare func()
		Expected []*domain.Monkey
		Err error
	} {
		{
			TestName: "list success",
			Prepare: func() {
				store.EXPECT().ListMonkeys().Return(
					[]*domain.Monkey {
						{
							Name: "knight",
							UUID: "123456",
						},
					},
					nil,
				)
			},
			Expected: []*domain.Monkey {
				{
					Name: "knight",
					UUID: "123456",
				},
			},
			Err: nil,
		},
		{
			TestName: "db error",
			Prepare: func() {
				store.EXPECT().ListMonkeys().Return(nil, er.New("db error"))
			},
			Expected: nil,
			Err: errors.NewCommon(domain.Source, er.New("db error")).SetCode(errors.ErrorCodeDatabaseError),
		},
	}

	for _, test := range testTable {
		test.Prepare()
		monkeys, err := controller.ListMonkeys()

		assert.EqualValues(test.Expected, monkeys)
		assert.EqualValues(test.Err, err)
	}
}


func TestRetrieveMonkey(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := getStore(ctrl)

	controller := control.NewMonkeyController(store, log.NewLogger("control_test"))

	type TestCase struct {
		TestName string
		Input string
		Prepare func(testCase TestCase)
		Expected *domain.Monkey
		Err error
	}

	validUUID := uuid.Must(uuid.NewV4()).String()

	testTable := []TestCase {
		{
			TestName: "retrive monkey success",
			Input: validUUID,
			Prepare: func(testCase TestCase) {
				store.EXPECT().RetrieveMonkey(testCase.Input).Return(
					&domain.Monkey{
						Name: "knight",
						UUID: validUUID,
					},
					nil,
				)
			},
			Expected: &domain.Monkey{
				Name: "knight",
				UUID: validUUID,
			},
			Err: nil,
		},
		{
			TestName: "db error",
			Input: validUUID,
			Prepare: func(testCase TestCase) {
				store.EXPECT().RetrieveMonkey(testCase.Input).Return(
					nil,
					er.New("db error"),
				)
			},
			Expected: nil,
			Err: errors.NewCommon(domain.Source, er.New("db error")).SetCode(errors.ErrorCodeDatabaseError),
		},
		{
			TestName: "invalid uuid",
			Input: "123456",
			Prepare: func(testCase TestCase) {},
			Expected: nil,
			Err: errors.New(domain.Source, errors.ErrorCodeInvalidArgs).AddFieldError("uuid", "uuid is invalid"),
		},
	}

	for _, test := range testTable {
		test.Prepare(test)
		monkey, err := controller.RetrieveMonkey(test.Input)

		assert.EqualValues(test.Expected, monkey)
		assert.EqualValues(test.Err, err)
	}
}

