package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/models"
	"github.com/turbex-backend/src/structs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func IncludeDatabaseConn(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(consts.CONTEXT_DB, database)
		c.Next()
	}
}

func IncludeGridFSBucket(bucket *gridfs.Bucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(consts.CONTEXT_GRIDFS, bucket)
		c.Next()
	}
}

func IncludeEnvironmentVariables(env *structs.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(consts.CONTEXT_ENV, env)
		c.Next()
	}
}

func IncludeSession(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(consts.SESSION_COOKIE_NAME)
		if err != nil {
			c.Next()
			return
		}
		userSession, found := getSession(database, cookie)
		if !found {
			c.Next()
			return
		}
		c.Set(consts.CONTEXT_SESSION, userSession)
		c.Next()
	}
}

func RequireLogged() gin.HandlerFunc {
	abortRequest := func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "You need to login frist"})
	}
	return func(c *gin.Context) {
		userSessionTmp, exist := c.Get(consts.CONTEXT_SESSION)
		if !exist {
			abortRequest(c)
			return
		}
		userSession, assertion := userSessionTmp.(*models.Session)
		if !assertion {
			abortRequest(c)
			return
		}
		expirationDate, err := time.Parse(consts.DATE_FORMAT, userSession.ExpirationDate)
		if err != nil || expirationDate.Before(time.Now()) {
			abortRequest(c)
			return
		}
		c.Next()
	}
}
