package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"go.mongodb.org/mongo-driver/mongo"
)

func IncludeDatabaseConn(database *mongo.Database) gin.HandlerFunc {
  return func(c *gin.Context) {
    c.Set(consts.CONTEXT_DB, database)
    c.Next()
  }
}
