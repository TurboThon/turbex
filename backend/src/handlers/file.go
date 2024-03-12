package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func DoUploadFile(c *gin.Context, db *mongo.Database, bucket *gridfs.Bucket, userSession *models.Session) {
	file, header, err := c.Request.FormFile("encrypted_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()
	if header.Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename cannot be empty"})
		return
	}
	encrypted_file_key := c.Request.FormValue("encrypted_file_key")
	if encrypted_file_key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'encrypted_file_key' is not empty"})
		return
	}

	objectID, err := bucket.UploadFromStream(header.Filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	fileShare := models.FileShare{
		Id:            primitive.NewObjectID(),
		UserRef:       userSession.UserId,
		FileRef:       objectID.Hex(),
		EncryptionKey: encrypted_file_key,
		// Expiration in one year: expiration is currently not supported by backend
		ExpirationDate: time.Now().Add(12 * 30 * 24 * time.Hour).UTC().Format(consts.DATE_FORMAT),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.Collection(consts.COLLECTION_FILE_SHARE).InsertOne(ctx, fileShare)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	c.String(http.StatusCreated, "Created")
}
