package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
    ID                 primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    Title              string               `bson:"title" json:"title"`
    Description        string               `bson:"description" json:"description"`
    Capacity           int                  `bson:"capacity" json:"capacity"`
    RegisteredUserIds  []primitive.ObjectID `bson:"registeredUserIds" json:"registeredUserIds"`
    CreatedByUserId    primitive.ObjectID   `bson:"createdByUserId" json:"createdByUserId"`
    CreatedAt          int64                `bson:"createdAt" json:"createdAt"`
}

type TaskView struct {
    ID               primitive.ObjectID `json:"id"`
    Title            string             `json:"title"`
    Description      string             `json:"description"`
    Capacity         int                `json:"capacity"`
    RegisteredCount  int                `json:"registeredCount"`
    Remaining        int                `json:"remaining"`
    Registrants      []UserPublic       `json:"registrants"`
    CreatedAt        int64              `json:"createdAt"`
}
