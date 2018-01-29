package bootstrap

import (
	"monkey/infra/diagnose"
	"monkey/infra/http"
	"monkey/infra/log"
	"monkey/inter"
)

func Server() {
	var (
		diagnoser *diagnose.HealthChecker
	)
	diagnoser, _ = diagnose.New()

	GetHTTPServer().
		AddEndpoint("/_diagnose", http.NewDiagnoser(diagnoser)).
		AddVersionEndpoint(
			1, "/monkey",
			inter.NewMonkeyHandler(nil, log.NewLogger("bootstrap")),
		).
		Start()
}
