package database

import (
	"context"
	"log"

	"github.com/afuradanime/backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to MongoDB
func InitMongoDB(Config config.Config) (client *mongo.Client, err error) {

	clientOptions := options.Client().ApplyURI(Config.MongoConnectionString)
	clientOptions.SetAuth(options.Credential{
		Username: Config.MongoUsername,
		Password: Config.MongoPassword,
	})

	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB successfully!")
	return client, nil
}
