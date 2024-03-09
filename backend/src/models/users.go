package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `json:"id,omitempty"`
	UserName   string             `json:"userName"`
	FirstName  string             `json:"firstName"`
	LastName   string             `json:"lastName"`
	Password   string             `json:"password"`
	PrivateKey string             `json:"privateKey"`
}

type APICreateUserRequest struct {
	UserName   string `json:"userName" validate:"required"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	Password   string `json:"password" validate:"required"`
	PrivateKey string `json:"privateKey" validate:"required"`
}

type APIUser struct {
 	UserName   string `json:"userName" validate:"required"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
}

type APIUserDetails struct {
	UserName   string `json:"userName" validate:"required"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	PrivateKey string `json:"privateKey" validate:"required"`
}
