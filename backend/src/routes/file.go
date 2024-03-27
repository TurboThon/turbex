package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/handlers"
	"github.com/turbex-backend/src/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// @BasePath /api/v1
// @Summary Encrypted file upload
// @Schemes
// @Description Uploads an encrypted file. The file should be encrypted using AES-GCM with a 256 bits key created for the file
// @Tags file
// @Param			encrypted_file	formData	file			true	"Encrypted file to upload"
// @Param     encrypted_file_key formData string true "Primary file key encrypted using the owner personnal key"
// @Param     ephemeral_pub_key formData string true "Ephemeral public key created to share with the owner"
// @Accept json
// @Produce json
// @Success 201 {string} TODO
// @Router /file [POST]
func uploadFileRoute(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	gridfs := c.MustGet(consts.CONTEXT_GRIDFS).(*gridfs.Bucket)
	userSession := c.MustGet(consts.CONTEXT_SESSION).(*models.Session)
	handlers.DoUploadFile(c, database, gridfs, userSession)
}

// @BasePath /api/v1
// @Summary List files
// @Schemes
// @Description List files accessible by current the user
// @Tags file
// @Accept json
// @Produce json
// @Success 200 {object} models.APISuccess[[]models.File]
// @Router /file [get]
func listFilesRoute(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	gridfs := c.MustGet(consts.CONTEXT_GRIDFS).(*gridfs.Bucket)
	userSession := c.MustGet(consts.CONTEXT_SESSION).(*models.Session)
	handlers.DoListFile(c, database, gridfs, userSession)
}
