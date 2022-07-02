package handler

import (
	"encoding/base64"

	"github.com/adiletelf/jwt-auth-go/internal/model"
)

type Handler struct {
	tr model.TokenRepo
}

func New(tr model.TokenRepo) *Handler {
	return &Handler{
		tr: tr,
	}
}

func encodeTokenBase64(td model.TokenDetails) model.TokenDetails {
	return model.TokenDetails{
		AccessToken:  base64.RawStdEncoding.EncodeToString([]byte(td.AccessToken)),
		RefreshToken: base64.RawStdEncoding.EncodeToString([]byte(td.RefreshToken)),
	}
}
