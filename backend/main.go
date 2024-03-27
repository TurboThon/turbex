package main

import (
	"context"
	"log"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/database"
	"github.com/turbex-backend/src/routes"
	"github.com/turbex-backend/src/routines"
	"github.com/turbex-backend/src/structs"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func main() {
	envVars := &structs.Env{}

	err := env.Parse(envVars)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Loaded environment variables")
	}

	client, err := database.InitDbConn(envVars)

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to mongodb")
	db := client.Database("turbex")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to gridfs bucket")

	database.SetupIndexes(db)

	go func() {
		for {
			routines.CleanExpiredSessions(db)
			time.Sleep(consts.SESSION_CLEAN_INTERVAL * time.Second)
		}
	}()
	log.Println("Session cleaning goroutine started")

	router := routes.SetupRouter(db, bucket, envVars)
	router.Run(":8000")
}
