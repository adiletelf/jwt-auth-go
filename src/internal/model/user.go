package model

import (
	"github.com/google/uuid"
)

type User struct {
	UUID         uuid.UUID `bson:"_id"`
	RefreshToken string
}
