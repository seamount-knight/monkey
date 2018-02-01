package bootstrap

import (
	"gopkg.in/doug-martin/goqu.v4"
	"monkey/control"
	"monkey/infra/log"
	"monkey/inter"
)

func GetController(db *goqu.Database) *control.MonkeyController {
	return control.NewMonkeyController(
		inter.NewMonkeyStore(db),
		log.NewLogger("inter"),
	)
}
