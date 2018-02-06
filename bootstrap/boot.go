package bootstrap

import (
	"github.com/spf13/viper"
	"gopkg.in/doug-martin/goqu.v4"
	"monkey/infra/diagnose"
	"monkey/infra/http"
	"monkey/infra/log"
	"monkey/inter"
	"time"
	"monkey/infra/db"
)

func Server() {
	var (
		logger     log.Logger
		diagnoser  *diagnose.HealthChecker
		dbclient   *goqu.Database
		controller inter.MonkeyController
	)
	dbclient, logger, controller = Bootstrap()

	diagnoser, _ = diagnose.New()
	diagnoser.Add(db.NewChecker(dbclient))

	logger.Infof("Starting %s...", viper.GetString("COMPONENT"))

	GetHTTPServer().
		AddEndpoint("/_diagnose", http.NewDiagnoser(diagnoser)).
		AddVersionEndpoint(
			1, "/monkeys",
			inter.NewMonkeyHandler(controller, log.NewLogger("inter")),
		).
		Start()
}

// Bootstrap bootstrap basic resources
func Bootstrap() (
	dbclient *goqu.Database,
	logger log.Logger,
	controller inter.MonkeyController) {
	var (
		boot    interface{}
		retries = 10
	)

	logger = log.NewLogger(viper.GetString("COMPONENT"))
	log.SetLevel(log.Leveldebug)

	// connect to the database
	boot = Retry(retries, "database", func() (interface{}, error) {
		return GetDB()
	}, logger)
	dbclient = boot.(*goqu.Database)

	controller = GetController(dbclient)

	return
}

func Retry(retries int, driverName string, driverBootstrap func() (interface{}, error), logger log.Logger) interface{} {
	var (
		err    error
		driver interface{}
		sleep  = time.Second * 2
	)
	for i := 0; i < retries; i++ {
		driver, err = driverBootstrap()
		if err == nil {
			return driver
		}
		logger.Errorf(
			"----- Connection to \"%s\" error: %s. Will sleep for %s",
			driverName, err, sleep,
		)
		time.Sleep(sleep)
	}
	logger.Errorf(
		"----- Fatal connection error to \"%s\" error: %s. Will quit now.",
		driverName, err,
	)
	panic(err)
}
