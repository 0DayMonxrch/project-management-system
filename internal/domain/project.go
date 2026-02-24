package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectMember struct {
	UserID bson.ObjectID `bson:"user_id" json:"user_id"`
	Role   Role               `bson:"role"    json:"role"`
}

type Project struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name"          json:"name"`
	Description string             `bson:"description"   json:"description"`
	Members     []ProjectMember    `bson:"members"       json:"members"`
	CreatedBy   bson.ObjectID `bson:"created_by"    json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at"    json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"    json:"updated_at"`
}