package repository

import (
	"context"
	"fmt"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"github.com/adiletelf/jwt-auth-go/internal/model"
	"github.com/adiletelf/jwt-auth-go/internal/token"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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

func (tr *TokenRepoImpl) Generate(uid uuid.UUID) (model.TokenDetails, error) {
	td, err := tr.generateTokenDetails(uid)
	if err != nil {
		return model.TokenDetails{}, err
	}

	err = tr.insert(uid, td.RefreshToken)
	if err != nil {
		return model.TokenDetails{}, err
	}

	return td, nil
}

func (tr *TokenRepoImpl) Refresh(td model.TokenDetails) (model.TokenDetails, error) {
	uid, err := tr.extractUUID(td)
	if err != nil {
		return model.TokenDetails{}, err
	}

	newToken, err := tr.generateTokenDetails(uid)
	if err != nil {
		return model.TokenDetails{}, err
	}

	err = tr.upsert(uid, newToken.RefreshToken)
	if err != nil {
		return model.TokenDetails{}, err
	}

	return model.TokenDetails{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	}, nil
}

func (tr *TokenRepoImpl) extractUUID(td model.TokenDetails) (uuid.UUID, error) {
	accessUUID, err := token.ExtractUUIDFromToken(td.AccessToken, tr.cfg.ApiSecret)
	if err != nil {
		return uuid.Nil, err
	}
	refreshUUID, err := token.ExtractUUIDFromToken(td.AccessToken, tr.cfg.ApiSecret)
	if err != nil {
		return uuid.Nil, err
	}

	if accessUUID != refreshUUID {
		return uuid.Nil, fmt.Errorf("uuid of tokens dont match")
	}

	uid, err := uuid.Parse(refreshUUID)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func (tr *TokenRepoImpl) generateTokenDetails(uid uuid.UUID) (model.TokenDetails, error) {
	accessToken, err := token.GenerateAccessToken(uid, tr.cfg)
	if err != nil {
		return model.TokenDetails{}, err
	}
	refreshToken, err := token.GenerateRefreshToken(uid, tr.cfg)
	if err != nil {
		return model.TokenDetails{}, err
	}

	return model.TokenDetails{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (tr *TokenRepoImpl) insert(uid uuid.UUID, refreshToken string) error {
	hashedToken, err := hashPassword(refreshToken)
	if err != nil {
		return err
	}
	insert := bson.M{"_id": uid, "refreshToken": hashedToken}
	_, err = tr.collection.InsertOne(tr.ctx, insert)
	if err != nil {
		return fmt.Errorf("error while inserting document with id: %v", uid)
	}
	return nil
}

func (tr *TokenRepoImpl) upsert(uid uuid.UUID, refreshToken string) error {
	hashedToken, err := hashPassword(refreshToken)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": uid}
	update := bson.M{"$set": bson.M{"refreshToken": hashedToken}}
	opts := options.Update().SetUpsert(true)
	_, err = tr.collection.UpdateOne(tr.ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("error while upserting document with id: %v", uid)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
