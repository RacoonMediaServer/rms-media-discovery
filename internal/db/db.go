package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	AccountDatabase
}

type database struct {
	cli      *mongo.Client
	db       *mongo.Database
	accounts *mongo.Collection
}

const databaseTimeout = 40 * time.Second
const databaseName = "rms-media-discovery"

func Connect(uri string) (Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect to db failed: %w", err)
	}

	if err = cli.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("connect to db failed: %w", err)
	}

	return &database{
		cli:      cli,
		db:       cli.Database(databaseName),
		accounts: cli.Database(databaseName).Collection("accounts"),
	}, nil
}
