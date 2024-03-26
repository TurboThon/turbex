package models

// Files are stored in a gridfs table
// File shares are stored in a mongoDB table

type FileShare struct {
	FileRef         string             `json:"fileRef"`
	UserName        string             `json:"userId"`
	EncryptionKey   string             `json:"encryptionKey"`
	EphemeralPubKey string             `json:"ephemeralPubKey"`
	ExpirationDate  string             `json:"expirationDate"`
	CanWrite        bool               `json:"canWrite"`
}

type APICreateFileShareRequest struct {
  // FileRef         string `json:"file_ref" validate:"required"`
  // UserRef         string `json:"userId" validate:"required"`
	EncryptionKey   string `json:"encryptionKey" validate:"required"`
	EphemeralPubKey string `json:"ephemeralPubKey" validate:"required"`
	// ExpirationDate  string `json:"expirationDate"`
	// CanWrite        bool   `json:"can_write"`
}
