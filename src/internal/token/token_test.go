package token

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateAccessToken(t *testing.T) {
	secret := "testsecret42"
	tokenMinuteLifespan := "15"
	testcases := []struct {
		uid uuid.UUID
	}{
		{uuid.New()},
		{uuid.New()},
		{uuid.New()},
	}

	for _, tc := range testcases {
		token, err := GenerateAccessToken(tc.uid, secret, tokenMinuteLifespan)
		if err != nil {
			t.Error(err)
		}
		if err = TokenValid(token, secret); err != nil {
			t.Error(err)
		}
	}

}

func TestGenerateRefreshToken(t *testing.T) {
	secret := "testsecret42"
	tokenHourLifespan := "24"

	testcases := []struct {
		uid uuid.UUID
	}{
		{uuid.New()},
		{uuid.New()},
		{uuid.New()},
	}

	for _, tc := range testcases {
		token, err := GenerateRefreshToken(tc.uid, secret, tokenHourLifespan)
		if err != nil {
			t.Error(err)
		}
		if err = TokenValid(token, secret); err != nil {
			t.Error(err)
		}
	}
}

func TestExtractUUIDFromToken(t *testing.T) {
	accessMinuteLifespan := "15"
	refreshHourLifespan := "24"
	secret := "testsecret42"
	testcases := []struct {
		uid uuid.UUID
	}{
		{uuid.New()},
		{uuid.New()},
		{uuid.New()},
	}
	for _, tc := range testcases {
		accessToken, _ := GenerateAccessToken(tc.uid, secret, accessMinuteLifespan)
		refreshToken, _ := GenerateRefreshToken(tc.uid, secret, refreshHourLifespan)

		accessUUID, err := ExtractUUIDFromToken(accessToken, secret)
		if err != nil {
			t.Error(err)
		}

		refreshUUID, err := ExtractUUIDFromToken(refreshToken, secret)
		if err != nil {
			t.Error(err)
		}

		if accessUUID != refreshUUID {
			t.Error("uuid of access, refresh tokens doesn't match")
		}
	}
}


func TestRequestTokenValid(t *testing.T) {
	tokenName := "accessToken"
	secret := "testsecret42"
	accessToken, _ := GenerateAccessToken(uuid.New(), secret, "15")

	req, _ := http.NewRequest("POST", "/generate", nil)
	q := req.URL.Query()
	q.Add(tokenName, accessToken)
	req.URL.RawQuery = q.Encode()

	err := RequestTokenValid(req, tokenName, secret)
	if err != nil {
		t.Error(err)
	}
}

func Test_extractToken(t *testing.T) {
	tokenName := "accessToken"
	secret := "testsecret42"
	accessToken, _ := GenerateAccessToken(uuid.New(), secret, "15")

	req, _ := http.NewRequest("POST", "/generate", nil)
	q := req.URL.Query()
	q.Add(tokenName, accessToken)
	req.URL.RawQuery = q.Encode()

	token := extractToken(req, tokenName)
	if token == "" {
		t.Error("couldn't extract token from request")
	}
}
