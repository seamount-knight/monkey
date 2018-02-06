package inter_test

import (
	"testing"
	"monkey/infra/httptest"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"monkey/mock"
	"monkey/inter"
	"monkey/infra/http"
	"monkey/infra/log"
	netHttp "net/http"
	"monkey/domain"

	"monkey/infra/errors"
)

func init() {
	log.SetLevel(log.Leveldebug)
	gin.SetMode(gin.ReleaseMode)
}

func getController(ctrl *gomock.Controller) (*mock.MockMonkeyController) {
	ctrlMock := mock.NewMockMonkeyController(ctrl)
	return ctrlMock
}

func getApp(controller inter.MonkeyController) *gin.Engine {
	return http.NewServer(
		http.Config{
			Host: "0.0.0.0",
			Port: "8080",
			Component: "test",
		},
		nil,
	).Init().
	AddVersionEndpoint(1, "/monkeys", inter.NewMonkeyHandler(controller, log.NewLogger("test"))).
	GetApp()
}

func TestHTTPListMonkeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctrlMock := getController(ctrl)
	app := getApp(ctrlMock)
	e := httptest.New(app, t)

	type TestCase struct {
		// test case name -- easier to find later
		TestName string
		// any inputs like uuids etc
		Input string
		// expected return from API
		ExpectedJSON interface{}
		// expected status from API
		Status int
		// implementation of mocking and API call
		Execute func(testCase TestCase)
	}

	testTable := []TestCase{
		{
			"List empty",
			"",
			[]interface{}{},
			netHttp.StatusOK,
			func(testCase TestCase) {

				ctrlMock.EXPECT().ListMonkeys().Return([]*domain.Monkey{}, nil)

				e.GET("/v1/monkeys").Expect().
					Status(testCase.Status).
					JSON().Equal(testCase.ExpectedJSON)

			},
		},
	}

	for _, testCase := range testTable {
		testCase.Execute(testCase)
	}
}

func TestHTTPRetrieveMonkey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCtrl := getController(ctrl)

	app := getApp(mockCtrl)

	e := httptest.New(app, t)

	type TestCase struct {
		// test case name -- easier to find later
		TestName string
		// any inputs like uuids etc
		Input string
		// expected return from API
		ExpectedJSON interface{}
		// expected status from API
		Status int
		// implementation of mocking and API call
		Execute func(testCase TestCase)
	}

	testTable := []TestCase {
		{
			TestName: "Retrieve Monkey",
			Input: "123",
			ExpectedJSON: map[string]interface{}{
				"errors": []error{errors.New(domain.Source, errors.ErrorCodeInvalidArgs)},
			},
			Status: netHttp.StatusBadRequest,
			Execute: func(testCase TestCase) {
				mockCtrl.EXPECT().RetrieveMonkey(testCase.Input).
				Return(nil, errors.New(domain.Source, errors.ErrorCodeInvalidArgs))

				e.GET("/v1/monkeys/"+ testCase.Input).Expect().
					Status(testCase.Status).
					JSON().Equal(testCase.ExpectedJSON)
			},
		},
	}

	for _, test := range testTable {
		test.Execute(test)
	}
}