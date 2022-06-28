package model

import (
	"github.com/google/uuid"
)

type TokenDetails struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenRepo interface {
	Generate(uuid.UUID) (TokenDetails, error)
	Refresh(TokenDetails) (TokenDetails, error)
}
