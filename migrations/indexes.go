package migrations

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func RunIndexes(db *mongo.Database, log *slog.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	indexes := []struct {
		collection string
		model      mongo.IndexModel
	}{
		// Users
		{
			collection: "users",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			collection: "users",
			model: mongo.IndexModel{
				Keys: bson.D{{Key: "verification_token", Value: 1}},
			},
		},
		{
			collection: "users",
			model: mongo.IndexModel{
				Keys: bson.D{{Key: "reset_token", Value: 1}},
			},
		},
		// Projects
		{
			collection: "projects",
			model: mongo.IndexModel{
				Keys: bson.D{{Key: "members.user_id", Value: 1}},
			},
		},
		// Tasks
		{
			collection: "tasks",
			model: mongo.IndexModel{
				Keys: bson.D{{Key: "project_id", Value: 1}},
			},
		},
		{
			collection: "tasks",
			model: mongo.IndexModel{
				Keys: bson.D{{Key: "assigned_to", Value: 1}},
			},
		},
		// Notes
		{
			collection: "notes",
			model: mongo.IndexModel{
				Keys: bson.D{{Key: "project_id", Value: 1}},
			},
		},
	}

	for _, idx := range indexes {
		_, err := db.Collection(idx.collection).Indexes().CreateOne(ctx, idx.model)
		if err != nil {
			return err
		}
		log.Info("index created", "collection", idx.collection, "keys", idx.model.Keys)
	}

	return nil
}