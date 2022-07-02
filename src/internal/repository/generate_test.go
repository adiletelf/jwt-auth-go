package repository

import (
	"context"
	"testing"

	"github.com/adiletelf/jwt-auth-go/internal/config"
	"github.com/adiletelf/jwt-auth-go/internal/token"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTokenRepoImpl_Generate(t *testing.T) {
	tr := NewTokenRepo(context.Background(), &config.Config{
		ApiSecret:                 "testsecret42",
		AccessTokenMinuteLifespan: "15",
		RefreshTokenHourLifespan:  "24",
	}, collectionStub{})

	testcases := []struct {
		uid uuid.UUID
	}{
		{uuid.New()},
		{uuid.New()},
		{uuid.New()},
	}

	for _, tc := range testcases {
		td, err := tr.Generate(tc.uid)
		if err != nil {
			t.Error(err)
		}

		assert.NotEqual(t, "", td.AccessToken)
		assert.NotEqual(t, "", td.RefreshToken)
		assert.Nil(t, token.TokenValid(td.AccessToken, tr.cfg.ApiSecret))
		assert.Nil(t, token.TokenValid(td.RefreshToken, tr.cfg.ApiSecret))
	}
}
