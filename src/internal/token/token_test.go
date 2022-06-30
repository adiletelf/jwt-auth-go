package token

import (
	"testing"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"github.com/google/uuid"
)

func TestGenerateAccessToken(t *testing.T) {
	cfg := &config.Config{
		AccessTokenMinuteLifespan: "15",
		ApiSecret:                 "testsecret42",
	}
	testcases := []struct {
		uid uuid.UUID
	}{
		{uuid.New()},
		{uuid.New()},
		{uuid.New()},
	}

	for _, tc := range testcases {
		token, err := GenerateAccessToken(tc.uid, cfg)
		if err != nil {
			t.Error(err)
		}
		if err = TokenValid(token, cfg.ApiSecret); err != nil {
			t.Error(err)
		}
	}

}

func TestGenerateRefreshToken(t *testing.T) {
	cfg := &config.Config{
		RefreshTokenHourLifespan: "24",
		ApiSecret:                "testsecret42",
	}
	testcases := []struct {
		uid uuid.UUID
	}{
		{uuid.New()},
		{uuid.New()},
		{uuid.New()},
	}

	for _, tc := range testcases {
		token, err := GenerateRefreshToken(tc.uid, cfg)
		if err != nil {
			t.Error(err)
		}
		if err = TokenValid(token, cfg.ApiSecret); err != nil {
			t.Error(err)
		}
	}
}

func TestExtractUUIDFromToken(t *testing.T) {
	cfg := &config.Config{
		AccessTokenMinuteLifespan: "15",
		RefreshTokenHourLifespan:  "24",
		ApiSecret:                 "testsecret42",
	}
	testcases := []struct {
		uid uuid.UUID
	}{
		{uuid.New()},
		{uuid.New()},
		{uuid.New()},
	}
	for _, tc := range testcases {
		accessToken, _ := GenerateAccessToken(tc.uid, cfg)
		refreshToken, _ := GenerateRefreshToken(tc.uid, cfg)

		accessUUID, err := ExtractUUIDFromToken(accessToken, cfg.ApiSecret)
		if err != nil {
			t.Error(err)
		}

		refreshUUID, err := ExtractUUIDFromToken(refreshToken, cfg.ApiSecret)
		if err != nil {
			t.Error(err)
		}

		if accessUUID != refreshUUID {
			t.Error("uuid of access, refresh tokens don't match")
		}
	}
}
