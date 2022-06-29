package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/adiletelf/jwt-auth-go/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	tr model.TokenRepo
}

func New(tr model.TokenRepo) *Handler {
	return &Handler{
		tr: tr,
	}
}

type GenerateQuery struct {
	UUID string `form:"uuid" json:"uuid" binding:"required"`
}

func (h *Handler) Generate(c *gin.Context) {
	id, err := parseUUID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenDetails, err := h.tr.Generate(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenDetails)
}

func parseUUID(c *gin.Context) (uuid.UUID, error) {
	var input GenerateQuery
	if err := c.ShouldBindQuery(&input); err != nil {
		return uuid.Nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(input.UUID)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(string(decoded))
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

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

	td := model.TokenDetails{
		AccessToken:  input.AccessToken,
		RefreshToken: input.RefreshToken,
	}
	newTokenDetails, err := h.tr.Refresh(td)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTokenDetails)
}
