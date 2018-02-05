package domain

import (
	"github.com/satori/go.uuid"
)

const TableName = "monkey"

type Monkey struct {
	UUID string	`json:"uuid" db:"uuid"`
	Name string	`json:"name" db:"name"`
}

func NewMonkey(name string) *Monkey {
	return &Monkey{
		Name: name,
	}
}

func (this *Monkey) GenerateUUID() *Monkey {
	this.SetUUID(uuid.Must(uuid.NewV4()).String())
	return this
}

func (this *Monkey) SetUUID(uuid string) *Monkey {
	this.UUID = uuid
	return this
}

type MonkeyStore interface {
	ListMonkeys() ([]*Monkey, error)
	RetrieveMonkey(uuid string) (*Monkey, error)
}
