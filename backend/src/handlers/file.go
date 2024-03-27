package handlers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "'encrypted_file_key' is empty"})
		return
	}
	ephemeral_pub_key := c.Request.FormValue("ephemeral_pub_key")
	if ephemeral_pub_key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'ephemeral_pub_key' is empty"})
		return
	}

	objectID, err := bucket.UploadFromStream(header.Filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	fileShare := models.FileShare{
		UserName:        userSession.UserName,
		FileRef:         objectID.Hex(),
		EncryptionKey:   encrypted_file_key,
		EphemeralPubKey: ephemeral_pub_key,
		// Expiration in one year: expiration is currently not supported by backend
		ExpirationDate: time.Now().Add(12 * 30 * 24 * time.Hour).UTC().Format(consts.DATE_FORMAT),
		CanWrite:       true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.Collection(consts.COLLECTION_FILE_SHARE).InsertOne(ctx, fileShare)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"fileid": objectID.Hex()})
}

func DoListFile(c *gin.Context, db *mongo.Database, bucket *gridfs.Bucket, userSession *models.Session) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := db.Collection(consts.COLLECTION_FILE_SHARE).Find(ctx, bson.M{"username": userSession.UserName})

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusOK, models.APISuccess[[]models.File]{Data: []models.File{}})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
		log.Println(err)
		return
	}
	defer cur.Close(ctx)

	var fileInfos []models.FileShare
	var fileIds []primitive.ObjectID

	for cur.Next(ctx) {
		var fileInfo models.FileShare
		err := cur.Decode(&fileInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
			log.Println(err)
			return
		}
		objectID, err := primitive.ObjectIDFromHex(fileInfo.FileRef)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
			log.Println(err)
			return
		}

		fileInfos = append(fileInfos, fileInfo)
		fileIds = append(fileIds, objectID)
	}

	// Convert to map for further use
	fileShareMap := map[string]models.FileShare{}
	for _, fileInfo := range fileInfos {
		fileShareMap[fileInfo.FileRef] = fileInfo
	}

	cur, err = bucket.Find(bson.M{"_id": bson.M{"$in": fileIds}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
		log.Println(err)
		return
	}
	defer cur.Close(ctx)

	filesMetadata := []models.APIFileMetadata{}
	// Check the user has correct permission over the file
	// We assume the document is present in the gridfs bucket
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = cur.All(ctx, &filesMetadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
		log.Println(err)
		return
	}

	// Populate other fileds of filesMetadata
	for index := range filesMetadata {
		fileInfo := fileShareMap[filesMetadata[index].ID]
		filesMetadata[index].CanWrite = fileInfo.CanWrite
		filesMetadata[index].EncryptionKey = fileInfo.EncryptionKey
		filesMetadata[index].EphemeralPubKey = fileInfo.EphemeralPubKey
		filesMetadata[index].ExpirationDate = fileInfo.ExpirationDate
	}

	c.JSON(http.StatusOK, models.APISuccess[[]models.APIFileMetadata]{Data: filesMetadata})
}

func DoGetFile(c *gin.Context, db *mongo.Database, bucket *gridfs.Bucket, userSession *models.Session) {
  // Check params
  fileID := c.Param("docid")
  // Check if the user has right to get the file
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
  
  var fileShare models.FileShare
  err := db.Collection(consts.COLLECTION_FILE_SHARE).FindOne(ctx, bson.M{"username": userSession.UserName, "fileref": fileID}).Decode(&fileShare)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusForbidden, models.APIError{Error: "You are not allowed to retrieve this document"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
		log.Println(err)
		return
	}

  objectID, err := primitive.ObjectIDFromHex(fileShare.FileRef)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
		log.Println(err)
		return
	}



  var buffer bytes.Buffer
  length, err := bucket.DownloadToStream(objectID, &buffer)
  if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
    log.Println(err)
    return
  }

  c.DataFromReader(http.StatusOK, length, "application/octet-stream", &buffer, map[string]string{})
}
