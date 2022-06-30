package token

import (
	"fmt"
	"strconv"
	"time"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateAccessToken(uid uuid.UUID, cfg *config.Config) (string, error) {
	tokenLifespan, err := strconv.Atoi(cfg.AccessTokenMinuteLifespan)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["uuid"] = uid
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(cfg.ApiSecret))
	return signedToken, err
}

func GenerateRefreshToken(uid uuid.UUID, cfg *config.Config) (string, error) {
	tokenLifespan, err := strconv.Atoi(cfg.RefreshTokenHourLifespan)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["uuid"] = uid
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(cfg.ApiSecret))
	return signedToken, err
}

func ExtractUUIDFromToken(token, secret string) (string, error) {
	tokenJwt, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := tokenJwt.Claims.(jwt.MapClaims)
	if !ok || !tokenJwt.Valid {
		return "", fmt.Errorf("invalid token")
	}

	uid, exists := claims["uuid"]
	if !exists {
		return "", fmt.Errorf("the field 'uuid' is not found")
	}
	return uid.(string), nil
}

func TokenValid(token, secret string) error {
	tokenJwt, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || tokenJwt == nil {
		return err
	}

	return nil
}