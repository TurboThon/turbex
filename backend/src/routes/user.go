package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

// @BasePath /api/v1
// @Summary List users
// @Schemes
// @Description Get a list of registered users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.APIUser
// @Router /user [get]
func listUsersRoute(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	handlers.DoListUsers(c, database)
}

// @BasePath /api/v1
// @Summary Create a user
// @Schemes
// @Description Register a new user
// @Tags user
// @Param request body models.APICreateUserRequest true "body"
// @Accept json
// @Produce json
// @Success 201 {string} Created
// @Router /user [POST]
func createUserRoute(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	handlers.DoAddUser(c, database)
}
