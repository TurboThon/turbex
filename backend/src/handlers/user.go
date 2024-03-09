package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DoListUsers(c *gin.Context, db *mongo.Database) {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  cur, err := db.Collection(consts.COLLECTION_USER).Find(ctx, bson.D{})

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    log.Print(err)
    return
  }
  defer cur.Close(ctx)

  var users []models.User

  for cur.Next(ctx) {
    var user models.User
    err := cur.Decode(&user)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      log.Print(err)
      return
    }
    users = append(users, user)
  }

  c.JSON(http.StatusOK, gin.H{})
}
