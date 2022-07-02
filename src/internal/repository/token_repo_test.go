package repository

import (
	"context"
	"testing"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionStub struct{}

func (cs collectionStub) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	var result mongo.InsertOneResult
	return &result, nil
}

func (cs collectionStub) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var result mongo.UpdateResult
	return &result, nil
}

func TestNewTokenRepo(t *testing.T) {
	stub := collectionStub{}
	got := NewTokenRepo(context.TODO(), &config.Config{
		ApiSecret:                 "testsecret42",
		AccessTokenMinuteLifespan: "15",
		RefreshTokenHourLifespan:  "24",
	}, stub)
	if got == nil {
		t.Error("couldn't initialize TokenRepoImpl")
	}
}
