package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Files are stored in a gridfs table
// This struct is used to extract files metadata
type File struct {
	Id         primitive.ObjectID    `json:"id" bson:"_id"`
	Length     int       `json:"length" bson:"length"`
	ChunkSize  int       `json:"chunkSize" bson:"chunkSize"`
	UploadDate time.Time `json:"uploadDate" bson:"uploadDate"`
	Filename   string    `json:"filename" bson:"filename"`
}

// Files are stored in a gridfs table
// File shares are stored in a mongoDB table
type FileShare struct {
	Id              primitive.ObjectID `json:"id" bson:"_id"`
	FileRef         string `json:"fileRef"`
	UserName        string `json:"userId"`
	EncryptionKey   string `json:"encryptionKey"`
	EphemeralPubKey string `json:"ephemeralPubKey"`
	ExpirationDate  string `json:"expirationDate"`
	CanWrite        bool   `json:"canWrite"`
}

type APIFileMetadata struct {
	Id              string    `json:"id" bson:"_id"`
	Length          int       `json:"length" bson:"length"`
	ChunkSize       int       `json:"chunkSize" bson:"chunkSize"`
	UploadDate      time.Time `json:"uploadDate" bson:"uploadDate"`
	Filename        string    `json:"filename" bson:"filename"`
	EncryptionKey   string    `json:"encryptionKey"`
	EphemeralPubKey string    `json:"ephemeralPubKey"`
	ExpirationDate  string    `json:"expirationDate"`
	CanWrite        bool      `json:"canWrite"`
}

type APICreateFileShareRequest struct {
	// FileRef         string `json:"file_ref" validate:"required"`
	// UserRef         string `json:"userId" validate:"required"`
	EncryptionKey   string `json:"encryptionKey" validate:"required"`
	EphemeralPubKey string `json:"ephemeralPubKey" validate:"required"`
	// ExpirationDate  string `json:"expirationDate"`
	// CanWrite        bool   `json:"can_write"`
}

type APIFileInfo struct {
	FileRef  string `json:"fileRef"`
	CanWrite bool   `json:"canWrite"`
}
