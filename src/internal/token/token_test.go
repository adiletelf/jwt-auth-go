package token

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateAccessToken(t *testing.T) {
	apiSecret := "testsecret42"
	tokenMinuteLifespan := "15"
	testcases := []struct {
		uid uuid.UUID
	}{
		{uuid.New()},
		{uuid.New()},
		{uuid.New()},
	}

	for _, tc := range testcases {
		token, err := GenerateAccessToken(tc.uid, tokenMinuteLifespan, apiSecret)
		if err != nil {
			t.Error(err)
		}
		if err = TokenValid(token, apiSecret); err != nil {
			t.Error(err)
		}
	}

}

func TestGenerateRefreshToken(t *testing.T) {
	apiSecret := "testsecret42"
	tokenHourLifespan := "24"

	testcases := []struct {
		uid uuid.UUID
	}{
		{uuid.New()},
		{uuid.New()},
		{uuid.New()},
	}

	for _, tc := range testcases {
		token, err := GenerateRefreshToken(tc.uid, tokenHourLifespan, apiSecret)
		if err != nil {
			t.Error(err)
		}
		if err = TokenValid(token, apiSecret); err != nil {
			t.Error(err)
		}
	}
}

func TestExtractUUIDFromToken(t *testing.T) {
	accessMinuteLifespan := "15"
	refreshHourLifespan := "24"
	apiSecret := "testsecret42"
	testcases := []struct {
		uid uuid.UUID
	}{
		{uuid.New()},
		{uuid.New()},
		{uuid.New()},
	}
	for _, tc := range testcases {
		accessToken, _ := GenerateAccessToken(tc.uid, accessMinuteLifespan, apiSecret)
		refreshToken, _ := GenerateRefreshToken(tc.uid, refreshHourLifespan, apiSecret)

		accessUUID, err := ExtractUUIDFromToken(accessToken, apiSecret)
		if err != nil {
			t.Error(err)
		}

		refreshUUID, err := ExtractUUIDFromToken(refreshToken, apiSecret)
		if err != nil {
			t.Error(err)
		}

		if accessUUID != refreshUUID {
			t.Error("uuid of access, refresh tokens don't match")
		}
	}
}

func Test_extractToken(t *testing.T) {
	accessToken, _ := GenerateAccessToken(uuid.New(), "15", "testsecret42")
	req, _ := http.NewRequest("GET", "/generate", nil)
	q := req.URL.Query()
	queryName := "accessToken"
	q.Add(queryName, accessToken)
	req.URL.RawQuery = q.Encode()

	token := extractToken(req, queryName)
	if token == "" {
		t.Error("couldn't extract token from request")
	}
}
