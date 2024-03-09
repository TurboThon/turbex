package main

import (
	"github.com/turbex-backend/src/routes"
)

func main() {
  router := routes.SetupRouter()
  router.Run(":8000")
}

