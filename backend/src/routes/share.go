package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/handlers"
	"github.com/turbex-backend/src/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// @BasePath /api/v1
// @Summary Share a document with a user
// @Schemes
// @Description Share a document with the specified user
// @Tags share
// @Param docid path string true "FileID"
// @Param username path string true "UserName"
// @Param request body models.APICreateFileShareRequest true "body"
// @Accept json
// @Produce json
// @Success 201 {string} Created
// @Router /share/{docid}/{username} [post]
func createShareRoute(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	userSession := c.MustGet(consts.CONTEXT_SESSION).(*models.Session)

	handlers.DoShareFile(c, database, userSession)
}
