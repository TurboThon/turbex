package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/handlers"
	"github.com/turbex-backend/src/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// @BasePath /api/v1
// @Summary List users
// @Schemes
// @Description Get a list of registered users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} []models.APIUserInfo
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

// @BasePath /api/v1
// @Summary Get a user by username
// @Schemes
// @Description Get all public information about a user
// @Tags user
// @Param username path string true "UserName"
// @Accept json
// @Produce json
// @Success 200 {object} models.APIUserPublic
// @Router /user/{username} [GET]
func getUserRoute(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	handlers.DoGetUserByUserName(c, database)
}


// @BasePath /api/v1
// @Summary Modify a user by username
// @Schemes
// @Description Change a user's properties, usefull to rotate keys or change password
// @Tags user
// @Param username path string true "UserName"
// @Param request body models.APIChangeUserRequest true "body"
// @Accept json
// @Produce json
// @Success 200 {string} Changed successfully
// @Router /user/{username} [PUT]
func modifyUserRoute(c *gin.Context) {
	database := c.MustGet(consts.CONTEXT_DB).(*mongo.Database)
	userSession := c.MustGet(consts.CONTEXT_SESSION).(*models.Session)

	handlers.DoChangeUser(c, database, userSession)
}
