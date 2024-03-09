package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/handlers"
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
