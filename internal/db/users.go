package db

import (
	"context"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

type UserDatabase interface {
	LoadUsers() ([]model.User, error)
	CreateUser(user model.User) error
	DeleteUser(id string) error
}

func (d *database) LoadUsers() (result []model.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	cur, err := d.users.Find(ctx, bson.D{})
	if err != nil {
		return
	}
	for cur.Next(ctx) {
		user := model.User{}
		if err = cur.Decode(&user); err != nil {
			return
		}
		result = append(result, user)
	}

	return
}

func (d *database) CreateUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	_, err := d.users.InsertOne(ctx, user)
	return err
}

func (d *database) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	_, err := d.users.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	return err
}
