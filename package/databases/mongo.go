package databases

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongoDB(user, password, host, port, dbName string) (*mongo.Database, error) {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin",
		user,
		password,
		host,
		port,
		dbName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return db, nil
}
