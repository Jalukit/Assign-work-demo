package db

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect(uri, dbName string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return err
    }
    if err := c.Ping(ctx, readpref.Primary()); err != nil {
        return err
    }

    Client = c
    DB = c.Database(dbName)

    // Indexes
    _ = ensureIndexes()
    return nil
}
