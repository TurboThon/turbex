package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// @BasePath /api/v1
// @Summary Purge database
// @Schemes
// @Description Delete all entries in the database, including files
// @Tags testing
// @Accept json
// @Produce json
// @Success 201 {string} List of errors
// @Router /admin/purgedb [DELETE]
func purgeDB(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	gridfs := c.MustGet(consts.CONTEXT_GRIDFS).(*gridfs.Bucket)

	handlers.DoPurgeDB(c, database, gridfs)
}

// @BasePath /api/v1
// @Summary Seed database
// @Schemes
// @Description Create testing accounts and files
// @Tags testing
// @Accept json
// @Produce json
// @Success 201 {string} List of errors
// @Router /admin/seeddb [POST]
func seedDB(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	gridfs := c.MustGet(consts.CONTEXT_GRIDFS).(*gridfs.Bucket)

	handlers.DoSeedDB(c, database, gridfs)
}
