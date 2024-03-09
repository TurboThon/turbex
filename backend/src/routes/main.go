package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
  swaggerfiles "github.com/swaggo/files"
	docs "github.com/turbex-backend/docs"
	"github.com/turbex-backend/src/handlers"
)

func notImplemented(c *gin.Context) {
  c.String(http.StatusNotImplemented, "API Endpoint not implemented yet")
}

func SetupRouter() *gin.Engine {
  r := gin.Default()

  docs.SwaggerInfo.BasePath = ""
  r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))


  apiV1 := r.Group("/api/v1")

  apiV1.GET("/health", handlers.DoStatus)

  // register a user
  apiV1.POST("/register", notImplemented)

  // List users
  apiV1.GET("/user", notImplemented)
  // Get a user by id
  apiV1.GET("/user/:id", notImplemented)
  // Modifies a user
  apiV1.PUT("/user/:id", notImplemented)

  // Returns a list of files
  apiV1.GET("/file", notImplemented)
  // Get a file by id
  apiV1.GET("/file/:id", notImplemented)
  // Uploads a file
  apiV1.POST("/file", notImplemented)
  // Delete a file
  apiV1.DELETE("/file/:id", notImplemented)

  // Share a file with a single user
  apiV1.POST("/share/:docid/:userid", notImplemented)
  // Share a file with multiple users
  apiV1.POST("/share/:docid", notImplemented)
  // Delete a share relation for a user to a file
  apiV1.DELETE("/share/:docid/:userid", notImplemented)

  admin := apiV1.Group("/admin")
  // Delete a user
  admin.DELETE("/user", notImplemented)
  admin.DELETE("/purgedb", notImplemented)
  admin.PUT("/seeddb", notImplemented)

  return r
}

