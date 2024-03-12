package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Files are stored in a gridfs table
// File shares are stored in a mongoDB table

type FileShare struct {
	Id             primitive.ObjectID `json:"id,omitempty"`
	FileRef        string             `json:"file_ref"`
	UserRef        string             `json:"userId"`
	EncryptionKey  string             `json:"encryptionKey"`
	ExpirationDate string             `json:"expirationDate"`
	CanWrite       bool               `json:"can_write"`
}
