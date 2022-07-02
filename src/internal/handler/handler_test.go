package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adiletelf/jwt-auth-go/internal/model"
	"github.com/adiletelf/jwt-auth-go/internal/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	generatePath string = "/generate"
	refreshPath  string = "/refresh"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	router = gin.Default()
	h := &Handler{
		tr: &TokenRepoStub{},
	}

	router.POST(generatePath, h.Generate)
	router.POST(refreshPath, h.Refresh)
	m.Run()
}

type TokenRepoStub struct{}

func (tr *TokenRepoStub) Generate(uid uuid.UUID) (model.TokenDetails, error) {
	return model.TokenDetails{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}, nil
}

func (tr *TokenRepoStub) Refresh(td model.TokenDetails) (model.TokenDetails, error) {
	return td, nil
}

func TestNew(t *testing.T) {
	stub := &TokenRepoStub{}
	got := New(stub)
	if got == nil {
		t.Error("couldn't initialize Handler")
	}
}

func TestHandler_Generate(t *testing.T) {
	req, _ := http.NewRequest("POST", generatePath, nil)
	q := req.URL.Query()
	q.Add("uuid", uuid.New().String())
	req.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body)
	validTokenDetails(t, response)
}

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

func validTokenDetails(t *testing.T, response map[string]any) {
	accessToken, atExists := response["accessToken"]
	refreshToken, rtExists := response["refreshToken"]

	assert.True(t, atExists)
	assert.True(t, rtExists)
	assert.NotEqual(t, "", accessToken)
	assert.NotEqual(t, "", refreshToken)
}