package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/adiletelf/jwt-auth-go/internal/model"
	"github.com/gin-gonic/gin"
)

type RefreshBody struct {
	AccessToken  string `json:"accessToken" binding:"required"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (h *Handler) Refresh(c *gin.Context) {
	var input RefreshBody
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	td, err := decodeRefreshBody(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newtoken, err := h.tr.Refresh(td)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, encodeTokenBase64(newtoken))
}

func decodeRefreshBody(input RefreshBody) (model.TokenDetails, error) {
	decodedAccessToken, err := base64.StdEncoding.DecodeString(input.AccessToken)
	if err != nil {
		return model.TokenDetails{}, err
	}
	decodedRefreshToken, err := base64.StdEncoding.DecodeString(input.RefreshToken)
	if err != nil {
		return model.TokenDetails{}, err
	}

	return model.TokenDetails{
		AccessToken:  string(decodedAccessToken),
		RefreshToken: string(decodedRefreshToken),
	}, nil
}
