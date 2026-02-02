package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserRole string

const (
    RoleUser  UserRole = "user"
    RoleAdmin UserRole = "admin"
)

type User struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Email        string             `bson:"email" json:"email"`
    PasswordHash string             `bson:"passwordHash" json:"-"`
    Role         UserRole           `bson:"role" json:"role"`
    CreatedAt    int64              `bson:"createdAt" json:"createdAt"`
}

type UserPublic struct {
    ID    primitive.ObjectID `json:"id"`
    Email string             `json:"email"`
    Role  UserRole           `json:"role"`
}

func (u User) Public() UserPublic {
    return UserPublic{ID: u.ID, Email: u.Email, Role: u.Role}
}
