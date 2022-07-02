package handler

import (
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

func validTokenDetails(t *testing.T, response map[string]any) {
	accessToken, atExists := response["accessToken"]
	refreshToken, rtExists := response["refreshToken"]

	assert.True(t, atExists)
	assert.True(t, rtExists)
	assert.NotEqual(t, "", accessToken)
	assert.NotEqual(t, "", refreshToken)
}
