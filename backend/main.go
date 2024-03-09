package main

import (
	"context"
	"log"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/turbex-backend/src/database"
	"github.com/turbex-backend/src/routes"
	"github.com/turbex-backend/src/structs"
)

func initEnv(env *structs.Env) {

}

func main() {
	envVars := &structs.Env{}

	err := env.Parse(envVars)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Loaded environment variables")
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

	log.Print("Connected to mongodb")

	router := routes.SetupRouter(client.Database("turbex"))
	router.Run(":8000")
}

