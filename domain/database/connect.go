package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(dburl string, dbname string) (*mongo.Database, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburl))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }

    return client.Database(dbname), nil
}
