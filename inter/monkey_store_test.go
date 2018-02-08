package inter_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gopkg.in/doug-martin/goqu.v4"
	"monkey/inter"
	"monkey/domain"
	"errors"
)

func TestListMonkeys(t *testing.T) {
	assert := assert.New(t)

	db, mock, _ := sqlmock.New()
	goquDB := goqu.New("postgres", db)
	store := inter.NewMonkeyStore(goquDB)

	defer db.Close()

	type TestCase struct {
		TestName string
		Prepare func()
		Expected []*domain.Monkey
		Err error
	}

	//querySQL := `SELECT "uuid", "name" FROM "monkey" ORDER BY "name" ASC`

	testTable := []TestCase {
		{
			TestName: "Successful case",
			Prepare: func() {
				rows := sqlmock.NewRows([]string{"name", "uuid"}).AddRow("knight", "123456")
				mock.ExpectQuery(".*").WillReturnRows(rows)
			},
			Expected: []*domain.Monkey{
				{
					Name: "knight",
					UUID: "123456",
				},
			},
			Err: nil,
		},
		{
			TestName: "Database error",
			Prepare: func() {
				mock.ExpectQuery(".*").WillReturnError(errors.New("database error"))
			},
			Expected: nil,
			Err: errors.New("database error"),
		},
	}

	for _, test := range testTable {
		test.Prepare()
		monkeys, err := store.ListMonkeys()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}

		assert.EqualValues(test.Expected, monkeys)
		assert.EqualValues(test.Err, err)
	}
}

func TestRetrieveMonkey(t *testing.T) {
	assert := assert.New(t)

	db, mock, _ := sqlmock.New()

	goquDB := goqu.New("postgres", db)

	store := inter.NewMonkeyStore(goquDB)

	type TestCase struct {
		TestName string
		Input string
		Prepare func(testCase TestCase)
		Expected *domain.Monkey
		Err error
	}

	querySQL := `.*`

	testCase := []TestCase {
		{
			TestName: "Successful case",
			Input: "123456",
			Prepare: func(testCase TestCase){
				rows := sqlmock.NewRows([]string{"name", "uuid"}).AddRow("knight", "123456")

				mock.ExpectQuery(querySQL).WithArgs(testCase.Input, 1).WillReturnRows(rows)
			},
			Expected: &domain.Monkey{
				Name: "knight",
				UUID: "123456",
			},
			Err: nil,
		},
		{
			TestName: "Database error",
			Input: "123456",
			Prepare: func(testCase TestCase) {
				mock.ExpectQuery(querySQL).WithArgs(testCase.Input, 1).WillReturnError(errors.New("database error"))
			},
			Expected: nil,
			Err: errors.New("database error"),
		},
	}

	for _, test := range testCase {
		test.Prepare(test)
		monkey, err := store.RetrieveMonkey(test.Input)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}

		assert.EqualValues(test.Expected, monkey)
		assert.EqualValues(test.Err, err)
	}
}