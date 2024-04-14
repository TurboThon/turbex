package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllUsers(db *mongo.Database) (*[]models.APIUserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := db.Collection(consts.COLLECTION_USER).Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []models.APIUserInfo

	for cur.Next(ctx) {
		var user models.APIUserInfo
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func GetUserDetailsByUserName(db *mongo.Database, username string) (*models.APIUserDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.APIUserDetails

	err := db.Collection(consts.COLLECTION_USER).FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserPublicByUserName(db *mongo.Database, username string) (*models.APIUserPublic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.APIUserPublic

	err := db.Collection(consts.COLLECTION_USER).FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DoListUsers(c *gin.Context, db *mongo.Database) {
	users, err := GetAllUsers(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Print(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": *users})
}

func DoAddUser(c *gin.Context, db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.APICreateUserRequest

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validator := validator.New()
	err = validator.Struct(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = GetUserPublicByUserName(db, user.UserName)
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusConflict, gin.H{"error": "A user with this username already exist"})
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}
	dbUser := models.User{
		Id:         primitive.NewObjectID(),
		UserName:   user.UserName,
		Password:   hashedPassword,
		PrivateKey: user.PrivateKey,
		PublicKey:  user.PublicKey,
	}

	_, err = db.Collection(consts.COLLECTION_USER).InsertOne(ctx, dbUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Print(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func DoGetUserByUserName(c *gin.Context, db *mongo.Database) {
	username := c.Param("username")

	user, err := GetUserPublicByUserName(db, username)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Print(err)
		return
	}
	c.JSON(http.StatusOK, *user)

}

func DoChangeUser(c *gin.Context, db *mongo.Database, session *models.Session) {
  // Get params
  userId := c.Param("id")

	var user models.APIChangeUserRequest

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validator := validator.New()
	err = validator.Struct(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

  if session.UserName != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to do that"})
		return
  }

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}
  if user.Password != "" {
    user.Password = hashedPassword
  }

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

  _, err = db.Collection(consts.COLLECTION_USER).UpdateOne(ctx, bson.M{"username": userId}, bson.M{"$set": user})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, models.APISuccess[string]{Data: "Resource changed successfully"})
}
