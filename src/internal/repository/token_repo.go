package repository

import (
	"context"

	"github.com/adiletelf/jwt-auth-go/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
)

type TokenRepoImpl struct {
	ctx        context.Context
	cfg        *config.Config
	collection *mongo.Collection
}

func NewTokenRepo(ctx context.Context, cfg *config.Config, collection *mongo.Collection) *TokenRepoImpl {
	return &TokenRepoImpl{
		ctx:        ctx,
		cfg:        cfg,
		collection: collection,
	}
}
