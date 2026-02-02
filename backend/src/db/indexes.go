package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// USERS: unique email
	users := DB.Collection("users")
	_, _ = users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetName("users_email_unique"),
	})

	// TASKS: sort by createdAt desc
	tasks := DB.Collection("tasks")
	_, _ = tasks.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "createdAt", Value: -1}},
		Options: options.Index().
			SetName("tasks_createdAt_desc"),
	})

	// TASKS: index for registeredUserIds
	_, _ = tasks.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "registeredUserIds", Value: 1}},
		Options: options.Index().
			SetBackground(true).
			SetName("registeredUserIds_idx"),
	})

	return nil
}
