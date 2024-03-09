package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/turbex-backend/docs"
	"github.com/turbex-backend/src/middlewares"
	"go.mongodb.org/mongo-driver/mongo"
)

func notImplemented(c *gin.Context) {
  c.String(http.StatusNotImplemented, "API Endpoint not implemented yet")
}

func setupDocs(r *gin.Engine) {
  docs.SwaggerInfo.BasePath = "/api/v1"
  r.GET("/api/swagger", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "/api/swagger/index.html")
  })
  r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func SetupRouter(database *mongo.Database) *gin.Engine {
  r := gin.Default()

  setupDocs(r)

  apiV1 := r.Group("/api/v1")

  apiV1.Use(middlewares.IncludeDatabaseConn(database))

  apiV1.GET("/health", healthRoute)

  // register a user
  apiV1.POST("/register", notImplemented)

  // List users
  apiV1.GET("/user", listUsersRoute)
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

