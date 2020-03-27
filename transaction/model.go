package transaction

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OperationType struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Description string             `json:"description" bson:"description" validate:"required"`
	IsCredit    bool               `json:"isCredit" bson:"isCredit"`
}

type Transaction struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	AccountID     primitive.ObjectID `json:"accountID" bson:"account" validate:"required"`
	OperationType primitive.ObjectID `json:"operationTypeID" bson:"operationType" validate:"required"`
	Amount        float64            `json:"amount" bson:"amount" validate:"required"`
	EventDate     string             `json:"eventDate,omitemptye" bson:"eventDate"`
}
