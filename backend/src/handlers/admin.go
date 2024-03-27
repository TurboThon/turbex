package handlers

import (
	"bytes"
	"context"
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

func DoPurgeDB(c *gin.Context, db *mongo.Database, bucket *gridfs.Bucket) {
  allErrors := []string{}

  collections := []string {
    consts.COLLECTION_FILE_SHARE,
    consts.COLLECTION_SESSIONS,
    consts.COLLECTION_USER,
  }

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
  for _, collection := range collections {
    _, err := db.Collection(collection).DeleteMany(ctx, bson.M{})
    if err != nil {
      allErrors = append(allErrors, err.Error())
    }
  }

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
  cur, err := bucket.GetFilesCollection().Find(ctx, bson.M{})
  if err != nil {
    allErrors = append(allErrors, err.Error())
  } else {
    allFiles := []models.File{}
    err = cur.All(ctx, &allFiles)
    if err != nil {
      allErrors = append(allErrors, err.Error())
    } else {
      for _, file := range allFiles {
        err = bucket.Delete(file.ID)
        if err != nil {
          allErrors = append(allErrors, err.Error())
        }
      }
    }
  }

  c.JSON(http.StatusNoContent, allErrors)
}

func DoSeedDB(c *gin.Context, db *mongo.Database, bucket *gridfs.Bucket) {
  allErrors := []string{}
	users := []interface{}{
		models.User{
      Id: primitive.NewObjectID(),
			FirstName:  "Bob",
			LastName:   "Lenon",
			UserName:   "bob",
			Password:   "pbkdf2_sha256:600000:mlLRD07RVGjGAXcaJp4BerlKbNTKEml5RAqnSBaQyRGCY7ds3RKu:4UZcy+Hu6d7OhaEfdtePT9rNSAu7VDGRiHz7xl5895F+FuR4Kl6NHLt+eT4Y4UPVJxdwTBO4K+eWPfn1o6IAIj3rRNMTfvPXxxGlvDRsrarkV8KfXTxHndfkwKEgruiDYibLA7dFK5qiQvVJDHl8jbvbMzjKKSZ4lX11o3SvQTo=", // Corresponding APIPassword: cGFzc3dvcmQ= (base64 of "password")
			PrivateKey: "nothing for now",
			PublicKey:  "nothing for now",
		},
		models.User{
      Id: primitive.NewObjectID(),
			FirstName:  "Alice",
			LastName:   "Bissi",
			UserName:   "alice",
			Password:   "pbkdf2_sha256:600000:1OQs5hgEoFtCnMX13kiC/AcOVnb6Qc4qdnkaZj9PxM0hFrw9x5G4:2fSBABbFbGq4OUL3vWqaRSSMZRbRM36WEy2J6BkCoS3qRSFgsIiwuys/2IF65W8zWJZAhhfySH/78BhjfouBtxBqLS2Sd673cNvbxNyw5yVvMVEJMpIGXTqwM1X1j4cuDO64wsMd4Bc9trgVah6ANXuI8W6/2YuBUYaN+jKT7/o=", // Corresponding APIPassword: cGFzc3dvcmQ= (base64 of "password")
			PrivateKey: "nothing for now",
			PublicKey:  "nothing for now",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

  _, err := db.Collection(consts.COLLECTION_USER).InsertMany(ctx, users)
  if err != nil {
    allErrors = append(allErrors, err.Error())
  }
  
  // Add a document for alice to bob
  buffer := bytes.NewBufferString("Content of txt file sent from alice to bob.")
  objectId, err := bucket.UploadFromStream("alice-to-bob.txt", buffer)
  if err != nil {
    allErrors = append(allErrors, err.Error())
  } else {
    fileShares := []interface{}{
      models.FileShare{
        ID: primitive.NewObjectID(),
        FileRef: objectId.Hex(),
        UserName: "alice",
        CanWrite: true,
        ExpirationDate: "nothing for now",
        EphemeralPubKey: "nothing for now",
        EncryptionKey: "nothing for now",
      },
      models.FileShare{
        ID: primitive.NewObjectID(),
        FileRef: objectId.Hex(),
        UserName: "bob",
        CanWrite: false,
        ExpirationDate: "nothing for now",
        EphemeralPubKey: "nothing for now",
        EncryptionKey: "nothing for now",
      },
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := db.Collection(consts.COLLECTION_FILE_SHARE).InsertMany(ctx, fileShares)
    if err != nil {
      allErrors = append(allErrors, err.Error())
    }
  }

  // Add document for alice only
  buffer = bytes.NewBufferString("Content of txt file sent for alice only.")
  objectId, err = bucket.UploadFromStream("alice.txt", buffer)
  if err != nil {
    allErrors = append(allErrors, err.Error())
  } else {
    fileShares := []interface{}{
      models.FileShare{
        ID: primitive.NewObjectID(),
        FileRef: objectId.Hex(),
        UserName: "alice",
        CanWrite: true,
        ExpirationDate: "nothing for now",
        EphemeralPubKey: "nothing for now",
        EncryptionKey: "nothing for now",
      },
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := db.Collection(consts.COLLECTION_FILE_SHARE).InsertMany(ctx, fileShares)
    if err != nil {
      allErrors = append(allErrors, err.Error())
    }
  }

	c.JSON(http.StatusCreated, allErrors)
}
