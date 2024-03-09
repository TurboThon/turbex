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

	go func() {
		for {
			routines.CleanExpiredSessions(db)
			time.Sleep(consts.SESSION_CLEAN_INTERVAL * time.Second)
		}
	}()
	log.Println("Session cleaning goroutine started")

	router := routes.SetupRouter(db, envVars)
	router.Run(":8000")
}
