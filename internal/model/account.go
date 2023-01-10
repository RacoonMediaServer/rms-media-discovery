package model

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Id          string `bson:"_id,omitempty"`
	Credentials map[string]string
	Limit       uint
}

func (a Account) Service() string {
	idx := strings.Index(a.Id, ".")
	if idx < 0 {
		return ""
	}
	return a.Id[:idx]
}

func (a Account) IsValid() bool {
	return strings.Contains(a.Id, ".")
}

func (a *Account) GenerateId(serviceId string) {
	a.Id = fmt.Sprintf("%s.%s", serviceId, uuid.NewV4().String())
}
