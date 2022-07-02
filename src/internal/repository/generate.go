package repository

import (
	"fmt"

	"github.com/adiletelf/jwt-auth-go/internal/model"
	"github.com/adiletelf/jwt-auth-go/internal/token"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (tr *TokenRepoImpl) Generate(uid uuid.UUID) (model.TokenDetails, error) {
	td, err := token.GenerateTokenDetails(uid, tr.cfg.ApiSecret, tr.cfg.AccessTokenMinuteLifespan, tr.cfg.RefreshTokenHourLifespan)
	if err != nil {
		return model.TokenDetails{}, err
	}

	err = tr.insert(uid, td.RefreshToken)
	if err != nil {
		return model.TokenDetails{}, err
	}

	return td, nil
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
