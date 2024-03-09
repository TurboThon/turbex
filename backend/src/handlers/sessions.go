package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/models"
	"github.com/turbex-backend/src/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func isValidLogin(db *mongo.Database, username string, password string) (*models.APIUserDetails, bool) {
	user, err := GetUserDetailsByUserName(db, username)
	log.Println(username, password)
	if err != nil {
		return nil, false
	}
	return user, true
}

func createSessionCookie(db *mongo.Database, userId string, duration time.Duration) (string, error) {
	randomBytes := make([]byte, 40)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(randomBytes)[:20]

	session := models.Session{
		Id:             primitive.NewObjectID(),
		CookieValue:    token,
		UserId:         userId,
		ExpirationDate: time.Now().Add(duration).UTC().Format(consts.DATE_FORMAT),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.Collection(consts.COLLECTION_SESSIONS).InsertOne(ctx, session)

	if err != nil {
		return "", err
	}

	return token, nil
}

func DoLogin(c *gin.Context, db *mongo.Database, env *structs.Env) {
	var loginData models.APILoginRequest

	err := c.BindJSON(&loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, isValid := isValidLogin(db, loginData.UserName, loginData.Password)
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong login data"})
		return
	}

	cookie, err := createSessionCookie(db, user.UserName, time.Duration(env.SessionDurationSeconds)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.SetCookie(consts.SESSION_COOKIE_NAME, cookie, env.SessionDurationSeconds, "/", c.GetHeader(consts.HOST_HEADER), true, true)
	c.JSON(http.StatusOK, user)
}
