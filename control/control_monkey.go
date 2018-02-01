package control

import (
	"monkey/domain"
	"monkey/infra/errors"
	"monkey/infra/log"
)

type MonkeyController struct {
	store domain.MonkeyStore
	log   log.Logger
}

func NewMonkeyController(store domain.MonkeyStore, log log.Logger) *MonkeyController {
	return &MonkeyController{
		store: store,
		log:   log,
	}
}

func (this *MonkeyController) ListMonkeys() ([]*domain.Monkey, error) {
	monkeys, err := this.store.ListMonkeys()
	if err != nil {
		return nil, errors.NewCommon(domain.Source, err).SetCode(errors.ErrorCodeDatabaseError)
	}
	return monkeys, nil
}

func (this *MonkeyController) RetrieveMonkey(uuid string) (*domain.Monkey, error) {
	monkey, err := this.RetrieveMonkey(uuid)
	if err != nil {
		return nil, errors.NewCommon(domain.Source, err).SetCode(errors.ErrorCodeDatabaseError)
	}
	return monkey, nil
}
