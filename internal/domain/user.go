package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role string

const (
	RoleAdmin        Role = "admin"
	RoleProjectAdmin Role = "project_admin"
	RoleMember       Role = "member"
)

type User struct {
	ID                bson.ObjectID `bson:"_id,omitempty"        json:"id"`
	Name              string             `bson:"name"                 json:"name"`
	Email             string             `bson:"email"                json:"email"`
	Password          string             `bson:"password"             json:"-"`
	Role              Role               `bson:"role"                 json:"role"`
	IsEmailVerified   bool               `bson:"is_email_verified"    json:"is_email_verified"`
	VerificationToken string             `bson:"verification_token"   json:"-"`
	ResetToken        string             `bson:"reset_token"          json:"-"`
	ResetTokenExpiry  time.Time          `bson:"reset_token_expiry"   json:"-"`
	RefreshToken      string             `bson:"refresh_token"        json:"-"`
	CreatedAt         time.Time          `bson:"created_at"           json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at"           json:"updated_at"`
}