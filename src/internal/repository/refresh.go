package repository

import (
	"fmt"

	"github.com/adiletelf/jwt-auth-go/internal/model"
	"github.com/adiletelf/jwt-auth-go/internal/token"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func (tr *TokenRepoImpl) Refresh(td model.TokenDetails) (model.TokenDetails, error) {
	uid, err := tr.extractUUID(td)
	if err != nil {
		return model.TokenDetails{}, err
	}

	newToken, err := token.GenerateTokenDetails(uid, tr.cfg.ApiSecret, tr.cfg.AccessTokenMinuteLifespan, tr.cfg.RefreshTokenHourLifespan)
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
