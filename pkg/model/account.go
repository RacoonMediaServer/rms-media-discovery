package model

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Id       string `gorm:"primaryKey"`
	Token    string
	Login    string
	Password string
	Limit    uint

	// Cookies - готовая сессия в формате заголовка Cookie ("a=b; c=d"), полученная в обычном браузере.
	// Используется вместо логина и пароля, когда сервис закрыт защитой от ботов (Cloudflare).
	Cookies string
	// UserAgent - User-Agent браузера, в котором получены Cookies (cf_clearance привязан к нему)
	UserAgent string
}

// HasSession сообщает, что аккаунт использует готовую сессию (cookies), а не логин и пароль
func (a *Account) HasSession() bool {
	return a.Cookies != ""
}

func (a *Account) Service() string {
	idx := strings.Index(a.Id, ".")
	if idx < 0 {
		return ""
	}
	return a.Id[:idx]
}

func (a *Account) IsValid() bool {
	return strings.Contains(a.Id, ".")
}

func (a *Account) GenerateId(serviceId string) {
	a.Id = fmt.Sprintf("%s.%s", serviceId, uuid.NewV4().String())
}

type Credentials struct {
	AccountId string
	Login     string
	Password  string

	// Cookies и UserAgent - готовая сессия (см. Account.Cookies). Если Cookies заданы, Login и Password могут быть пустыми.
	Cookies   string
	UserAgent string
}

type ApiKey struct {
	AccountId string
	Key       string
}

type AccessProvider interface {
	GetCredentials(serviceId string) (Credentials, error)
	GetApiKey(serviceId string) (ApiKey, error)
	MarkUnaccesible(accountId string)
}
