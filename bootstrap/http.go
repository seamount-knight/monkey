package bootstrap

import (
	"github.com/spf13/viper"
	"monkey/infra/http"

	"monkey/infra/log"
)

func getConfig() http.Config {
	return http.Config{
		Host:           viper.GetString("HOST"),
		Port:           viper.GetString("PORT"),
		AddHealthCheck: true,
		Component:      viper.GetString("COMPONENT"),
	}
}

func GetHTTPServer() *http.Server {
	return http.NewServer(
		getConfig(),
		log.NewLogger("bootstrap"),
	).Init()
}
