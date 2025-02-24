package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI    string
	DBName string
}

func GetMongoDatabase(ctx context.Context, config MongoConfig) (*mongo.Database, *mongo.Client, error) {
	if config.URI == "" {
		return nil, nil, fmt.Errorf("You must set your 'MONGODB_URI' environment variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		return nil, nil, err
	}

	dbNames, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		client.Disconnect(ctx)
		return nil, nil, err
	}

	exists := false
	for _, name := range dbNames {
		if name == config.DBName {
			exists = true
			break
		}
	}
	if !exists {
		client.Disconnect(ctx)
		return nil, nil, fmt.Errorf("Mongo. Database %s does not exist", config.DBName)
	}

	db := client.Database(config.DBName)
	return db, client, nil
}
