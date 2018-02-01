package inter

import (
	"gopkg.in/doug-martin/goqu.v4"
	"monkey/domain"
)

type MonkeyStore struct {
	goquDB *goqu.Database
}

func NewMonkeyStore(goquDB *goqu.Database) *MonkeyStore {
	return &MonkeyStore{
		goquDB: goquDB,
	}
}

func (this *MonkeyStore) ListMonkeys() ([]*domain.Monkey, error) {
	var monkeys []*domain.Monkey
	queryset := this.goquDB.From(domain.TableName).
		Select(
			goqu.I("uuid"),
			goqu.I("name"),
		).Order(goqu.I("name").Asc())

	if err := queryset.ScanStructs(&monkeys); err != nil {
		return nil, err
	}
	return monkeys, nil
}

func (this *MonkeyStore) RetrieveMonkey(uuid string) (*domain.Monkey, error) {
	var monkey domain.Monkey
	queryset := this.goquDB.From(domain.TableName).
		Select(
			goqu.I("uuid"),
			goqu.I("name"),
		).
		Order(
			goqu.I("name").Asc(),
		).Where(
		goqu.Ex{
			"uuid": uuid,
		},
	)
	if _, err := queryset.ScanStruct(&monkey); err != nil {
		return nil, err
	}
	return &monkey, nil
}
