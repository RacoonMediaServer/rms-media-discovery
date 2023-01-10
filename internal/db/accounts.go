package db

import (
	"context"

	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

type AccountDatabase interface {
	LoadAccounts() ([]model.Account, error)
	CreateAccount(account model.Account) error
	DeleteAccount(id string) error
}

func (d *database) LoadAccounts() (result []model.Account, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	cur, err := d.accounts.Find(ctx, bson.D{})
	if err != nil {
		return
	}
	for cur.Next(ctx) {
		user := model.Account{}
		if err = cur.Decode(&user); err != nil {
			return
		}
		result = append(result, user)
	}

	return
}

func (d *database) CreateAccount(user model.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	_, err := d.accounts.InsertOne(ctx, user)
	return err
}

func (d *database) DeleteAccount(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	_, err := d.accounts.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	return err
}
