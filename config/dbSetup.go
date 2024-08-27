package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func SetupDatabase() (*mongo.Database, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://uoozer_user:SLaOHytNuziKnJ8b@cluster0.dmg6t.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerSelectionTimeout(5*time.
		Second))
	if err != nil {
		cancel()
		return nil, nil, err
	}
	db := client.Database(EnvConfigs.DbName)

	fmt.Println("database instance is ready")

	// Create a TTL index on the "expiry" field of the "documents" collection
	collection := db.Collection("sms_list")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "expiry", Value: 1},
		},
		Options: options.Index().SetExpireAfterSeconds(0), // TTL of 3600 seconds (1 hour)
	}

	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	fmt.Println("TTL index created successfully")

	return db, cancel, nil
}
