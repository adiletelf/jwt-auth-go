package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"github.com/adiletelf/jwt-auth-go/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (tr *TokenRepoImpl) Generate(uuid uuid.UUID) (model.TokenDetails, error) {
	accessToken, err := generateAccessToken(uuid, tr.cfg)
	if err != nil {
		return model.TokenDetails{}, err
	}
	refreshToken, err := generateRefreshToken(uuid, tr.cfg)
	if err != nil {
		return model.TokenDetails{}, err
	}

	_, err = tr.upsert(uuid, refreshToken)
	if err != nil {
		return model.TokenDetails{}, err
	}

	return model.TokenDetails{
		// AccessToken:  base64.StdEncoding.EncodeToString([]byte(accessToken)),
		// RefreshToken: base64.StdEncoding.EncodeToString([]byte(refreshToken)),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (tr *TokenRepoImpl) Refresh(td model.TokenDetails) (model.TokenDetails, error) {
	return model.TokenDetails{
		AccessToken:  "newAccessToken",
		RefreshToken: "newRefreshToken",
	}, nil
}

func generateAccessToken(uuid uuid.UUID, cfg *config.Config) (string, error) {
	tokenLifespan, err := strconv.Atoi(cfg.AccessTokenMinuteLifespan)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["uuid"] = uuid
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(cfg.ApiSecret))
	return signedToken, err
}

func generateRefreshToken(uuid uuid.UUID, cfg *config.Config) (string, error) {
	tokenLifespan, err := strconv.Atoi(cfg.RefreshTokenHourLifespan)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["uuid"] = uuid
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(cfg.ApiSecret))
	return signedToken, err
}

func (tr *TokenRepoImpl) upsert(uid uuid.UUID, refreshToken string) (string, error) {
	hashedToken, err := hashPassword(refreshToken)
	if err != nil {
		return "", err
	}

	filter := bson.M{"_id": uid}
	update := bson.M{"$set": bson.M{"refreshToken": hashedToken}}
	opts := options.Update().SetUpsert(true)
	result, err := tr.collection.UpdateOne(tr.ctx, filter, update, opts)
	if err != nil {
		return "", err
	}

	upsertedID, ok := result.UpsertedID.(primitive.Binary)
	if !ok {
		return "", fmt.Errorf("error while upserting user with uuid: %v", uid)
	}

	return string(upsertedID.Data), nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
