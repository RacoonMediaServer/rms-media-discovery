package db

import "git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"

type UserDatabase interface {
	LoadUsers() ([]model.User, error)
	CreateUser(user model.User) error
	DeleteUser(id string) error
}
