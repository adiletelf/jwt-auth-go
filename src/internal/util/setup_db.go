package util

import (
	"context"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCollection(ctx context.Context, cfg *config.Config) (*mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI(cfg.DB.ConnectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	collection := client.Database(cfg.DB.Name).Collection(cfg.DB.CollectionName)
	return collection, nil
}
