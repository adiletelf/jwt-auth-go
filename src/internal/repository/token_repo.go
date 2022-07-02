package repository

import (
	"context"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"golang.org/x/crypto/bcrypt"

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

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
