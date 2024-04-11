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

func DoShareFile(c *gin.Context, db *mongo.Database, userSession *models.Session) {
	// Get Params
	userName := c.Param("username")
	docId := c.Param("docid")

	var createFileShare models.APICreateFileShareRequest

	err := c.BindJSON(&createFileShare)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validator := validator.New()
	err = validator.Struct(&createFileShare)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check the user has correct permission over the file
	// We assume the document is present in the gridfs bucket
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var currentUserPermissions models.FileShare

	err = db.Collection(consts.COLLECTION_FILE_SHARE).FindOne(ctx, bson.M{"fileref": docId, "username": userSession.UserName}).Decode(&currentUserPermissions)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The document does not exist or you don't have any permissions over it"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}
	if !currentUserPermissions.CanWrite {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to share this document"})
		return
	}

	// TODO: check the target user exist

	// Insert the share in the database
	fileShare := models.FileShare{
        Id:              primitive.NewObjectID(),
		UserName:        userName,
		FileRef:         docId,
		EncryptionKey:   createFileShare.EncryptionKey,
		EphemeralPubKey: createFileShare.EphemeralPubKey,
		ExpirationDate:  currentUserPermissions.ExpirationDate,
		CanWrite:        false,
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.Collection(consts.COLLECTION_FILE_SHARE).InsertOne(ctx, fileShare)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	c.String(http.StatusCreated, "Created")

}
