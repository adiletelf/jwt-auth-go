package repository

import (
	"context"
	"testing"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"github.com/adiletelf/jwt-auth-go/internal/model"
	"github.com/adiletelf/jwt-auth-go/internal/token"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTokenRepoImpl_Refresh(t *testing.T) {
	tr := NewTokenRepo(context.Background(), &config.Config{
		ApiSecret:                 "testsecret42",
		AccessTokenMinuteLifespan: "15",
		RefreshTokenHourLifespan:  "24",
	}, collectionStub{})

	testcases := []struct {
		td model.TokenDetails
	}{
		{getTokenDetails(t, tr.cfg.ApiSecret, tr.cfg.AccessTokenMinuteLifespan, tr.cfg.RefreshTokenHourLifespan)},
		{getTokenDetails(t, tr.cfg.ApiSecret, tr.cfg.AccessTokenMinuteLifespan, tr.cfg.RefreshTokenHourLifespan)},
		{getTokenDetails(t, tr.cfg.ApiSecret, tr.cfg.AccessTokenMinuteLifespan, tr.cfg.RefreshTokenHourLifespan)},
	}

	for _, tc := range testcases {
		newToken, err := tr.Refresh(tc.td)
		if err != nil {
			t.Error(err)
		}
		assert.Nil(t, token.TokenValid(newToken.AccessToken, tr.cfg.ApiSecret))
		assert.Nil(t, token.TokenValid(newToken.RefreshToken, tr.cfg.ApiSecret))
	}
}

func getTokenDetails(t *testing.T, secret, accessMinuteLifespan, refreshHourLifespan string) model.TokenDetails {
	td, err := token.GenerateTokenDetails(uuid.New(), secret, accessMinuteLifespan, refreshHourLifespan)
	if err != nil {
		t.Fatal(err)
	}
	return model.TokenDetails{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}
}
