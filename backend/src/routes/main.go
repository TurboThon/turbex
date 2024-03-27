package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/turbex-backend/docs"
	"github.com/turbex-backend/src/middlewares"
	"github.com/turbex-backend/src/structs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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

func SetupRouter(database *mongo.Database, bucket *gridfs.Bucket, env *structs.Env) *gin.Engine {
	r := gin.Default()

	setupDocs(r)

	apiV1 := r.Group("/api/v1")

	apiV1.Use(middlewares.IncludeDatabaseConn(database), middlewares.IncludeEnvironmentVariables(env), middlewares.IncludeSession(database))

	apiV1.GET("/health", healthRoute)

	apiV1.GET("/me", middlewares.RequireLogged(), meRoute)
	apiV1.POST("/login", loginRoute)
	// register a user
	apiV1.POST("/user", createUserRoute)

	// List users
	apiV1.GET("/user", middlewares.RequireLogged(), listUsersRoute)
	// Get a user by id
	apiV1.GET("/user/:username", middlewares.RequireLogged(), getUserRoute)
	// Modifies a user
	apiV1.PUT("/user/:id", middlewares.RequireLogged(), notImplemented)

	// Returns a list of files
	apiV1.GET("/file", middlewares.IncludeGridFSBucket(bucket), middlewares.RequireLogged(), listFilesRoute)
	// Uploads an encrypted file
	apiV1.POST("/file", middlewares.IncludeGridFSBucket(bucket), middlewares.RequireLogged(), uploadFileRoute)
	// Get a file by id
	apiV1.GET("/file/:id", middlewares.IncludeGridFSBucket(bucket), middlewares.RequireLogged(), notImplemented)
	// Delete a file
	apiV1.DELETE("/file/:id", middlewares.RequireLogged(), notImplemented)

	// Share a file with a single user
	apiV1.POST("/share/:docid/:username", middlewares.RequireLogged(), createShareRoute)
	// Share a file with multiple users
	apiV1.POST("/share/:docid", middlewares.RequireLogged(), notImplemented)
	// Delete a share relation for a user to a file
	apiV1.DELETE("/share/:docid/:username", middlewares.RequireLogged(), notImplemented)

	admin := apiV1.Group("/admin", middlewares.RequireLogged())
	// Delete a user
	admin.DELETE("/user", notImplemented)
	admin.DELETE("/purgedb", notImplemented)
	admin.PUT("/seeddb", notImplemented)

	return r
}
