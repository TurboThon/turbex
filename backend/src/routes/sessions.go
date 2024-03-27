package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/handlers"
	"github.com/turbex-backend/src/models"
	"github.com/turbex-backend/src/structs"
	"go.mongodb.org/mongo-driver/mongo"
)

// @BasePath /api/v1
// @Summary Login
// @Schemes
// @Description Log into the system using username and password
// @Tags login
// @Param request body models.APILoginRequest true "body"
// @Accept json
// @Produce json
// @Success 200 {object} models.APIUserDetails
// @Router /login [POST]
func loginRoute(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	env := c.MustGet(consts.CONTEXT_ENV).(*structs.Env)
	handlers.DoLogin(c, database, env)
}

// @BasePath /api/v1
// @Summary Get current user details
// @Schemes
// @Description Returns the current user details including keys
// @Tags login
// @Accept json
// @Produce json
// @Success 200 {object} models.APISuccess[models.APIUserDetails]
// @Router /me [GET]
func meRoute(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	userSession := c.MustGet(consts.CONTEXT_SESSION).(*models.Session)
	handlers.GetCurrentUser(c, database, userSession)
}
