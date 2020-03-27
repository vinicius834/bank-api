package account

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID             primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	DocumentNumber string             `json:"documentNumber" bson:"documentNumber" validate:"required"`
}
