package account

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID             primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Document       string             `json:"document" bson:"document" validate:"required"`
	AvalaibleLimit float64            `json:"avalaible_limit" bson:"avalaible_limit" validate:"required"`
}
