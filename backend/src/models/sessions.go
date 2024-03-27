package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	CookieValue    string             `json:"cookieValue"`
	UserName       string             `json:"userName"`
	ExpirationDate string             `json:"expirationDate"`
}

type APILoginRequest struct {
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
}
