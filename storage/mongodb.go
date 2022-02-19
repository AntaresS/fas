package storage

import (
	"context"
	"fas/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// NewClient creates a new mongodb conn and returns the client
func NewClient(db config.MongoDB) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.Uri))
	if err!=nil {
		return nil, err
	}

	return client, nil
}