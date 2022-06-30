package token

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateAccessToken(uid uuid.UUID, tokenMinuteLifespan, apiSecret string) (string, error) {
	tokenLifespan, err := strconv.Atoi(tokenMinuteLifespan)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["uuid"] = uid
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(apiSecret))
	return signedToken, err
}

func GenerateRefreshToken(uid uuid.UUID, tokenHourLifespan, apiSecret string) (string, error) {
	tokenLifespan, err := strconv.Atoi(tokenHourLifespan)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["uuid"] = uid
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(apiSecret))
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

func RequestTokenValid(r *http.Request, tokenName, secret string) error {
	token := extractToken(r, tokenName)
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func extractToken(r *http.Request, tokenName string) string {
	token := r.URL.Query().Get(tokenName)
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	splitted := strings.Split(bearerToken, " ")
	if len(splitted) == 2 {
		return splitted[1]
	}
	return ""
}
