package database

import (
	"context"
	"fmt"
	"time"

	"github.com/turbex-backend/src/structs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDbConn(env *structs.Env) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client()
	clientOptions.SetHosts([]string{fmt.Sprintf("%s:%d", env.DBHost, env.DBPort)})

	// If authentication is empty, use anonymous binding, the default when a mongo
	// instance is created
	if env.DBUser != "" && env.DBPass != "" {
		clientOptions.SetAuth(options.Credential{
			Username: env.DBUser,
			Password: env.DBPass,
		})
	}
	client, err := mongo.Connect(ctx, clientOptions)

	return client, err
}
