package model

type User struct {
	Id      string `bson:"_id,omitempty"`
	Info    string
	IsAdmin bool
}
