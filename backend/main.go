package main

import (
	"log"

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
  } else {
    log.Print("Connection to mongo db successful")
  }

  router := routes.SetupRouter(client.Database("turbex"))
  router.Run(":8000")
}

