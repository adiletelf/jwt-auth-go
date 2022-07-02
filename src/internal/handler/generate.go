package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

	c.JSON(http.StatusOK, encodeTokenBase64(tokenDetails))
}

func parseUUID(c *gin.Context) (uuid.UUID, error) {
	var input GenerateQuery
	if err := c.ShouldBindQuery(&input); err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(input.UUID)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
