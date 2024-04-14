package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupIndexes(db *mongo.Database) {
	usersIndexes := db.Collection(consts.COLLECTION_USER).Indexes()
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "username", Value: 1}},
    Options: options.Index().SetUnique(true),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := usersIndexes.CreateOne(ctx, indexModel)
	if err != nil {
		log.Println(err)
	}
}

func InitDbConn(env *structs.Env) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client()
	clientOptions.SetHosts([]string{fmt.Sprintf("%s:%d", env.DBHost, env.DBPort)})
  clientOptions.SetAppName("Turbex-backend")

	// If authentication is empty, use anonymous binding, the default when a mongo
	// instance is created
	if env.DBUser != "" && env.DBPass != "" {
		clientOptions.SetAuth(options.Credential{
      AuthSource: env.DBName,
      AuthMechanism: "SCRAM-SHA-256",
			Username: env.DBUser,
			Password: env.DBPass,
		})
	}
	client, err := mongo.Connect(ctx, clientOptions)

	return client, err
}
