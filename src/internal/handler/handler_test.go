package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adiletelf/jwt-auth-go/internal/model"
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

	router.GET(generatePath, h.Generate)
	router.GET(refreshPath, h.Refresh)
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
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", generatePath, nil)
	q := req.URL.Query()
	q.Add("uuid", uuid.New().String())
	req.URL.RawQuery = q.Encode()
	router.ServeHTTP(w, req)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	accessToken, atExists := response["accessToken"]
	refreshToken, rtExists := response["refreshToken"]

	assert.Nil(t, err)
	assert.True(t, atExists)
	assert.True(t, rtExists)
	assert.Equal(t, http.StatusOK, w.Code)
	if assert.NotEmpty(t, w.Body) {
		assert.NotEqual(t, "", accessToken)
		assert.NotEqual(t, "", refreshToken)
	}
}
