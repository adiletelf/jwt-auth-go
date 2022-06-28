package handler

import (
	"context"

	"github.com/adiletelf/jwt-auth-go/internal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	ctx        context.Context
	collection *mongo.Collection
}

func New(ctx context.Context, collection *mongo.Collection) *Handler {
	return &Handler{
		ctx:        ctx,
		collection: collection,
	}
}

func (h *Handler) Generate(uuid uuid.UUID) (model.TokenDetails, error) {
	return model.TokenDetails{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}, nil
}

func (h *Handler) Refresh(td model.TokenDetails) (model.TokenDetails, error) {
	return td, nil
}
