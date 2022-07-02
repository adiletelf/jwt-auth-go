package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adiletelf/jwt-auth-go/internal/token"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Refresh(t *testing.T) {
	input := getRefreshBody(t)
	jsonData, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", refreshPath, bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	validTokenDetails(t, response)
}

func getRefreshBody(t *testing.T) RefreshBody {
	td, err := token.GenerateTokenDetails(uuid.New(), "testsecret", "15", "24")
	if err != nil {
		t.Fatal(err)
	}
	encoded := encodeTokenBase64(td)
	return RefreshBody{
		AccessToken:  encoded.AccessToken,
		RefreshToken: encoded.RefreshToken,
	}
}
