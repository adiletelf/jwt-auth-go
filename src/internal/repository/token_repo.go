package repository

import (
	"context"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type TokenRepoImpl struct {
	ctx        context.Context
	cfg        *config.Config
	collection collection
}

func NewTokenRepo(ctx context.Context, cfg *config.Config, collection collection) *TokenRepoImpl {
	return &TokenRepoImpl{
		ctx:        ctx,
		cfg:        cfg,
		collection: collection,
	}
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
