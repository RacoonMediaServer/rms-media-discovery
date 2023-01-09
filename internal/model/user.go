package model

import "time"

type User struct {
	Id          string
	Info        string
	LastRequest time.Time
}
