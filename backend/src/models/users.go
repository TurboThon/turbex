package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	UserName   string             `json:"userName"`
	Password   string             `json:"password"`
	PrivateKey string             `json:"privateKey"`
	PublicKey  string             `json:"publicKey"`
}

type APICreateUserRequest struct {
	UserName   string `json:"userName" validate:"required"`
	Password   string `json:"password" validate:"required"`
	PrivateKey string `json:"privateKey" validate:"required"`
	PublicKey  string `json:"publicKey" validate:"required"`
}

type APIChangeUserRequest struct {
	Password   string `json:"password" bson:"password,omitempty"`
	PrivateKey string `json:"privateKey" bson:"privatekey,omitempty"`
	PublicKey  string `json:"publicKey" bson:"publickey,omitempty"`
}

type APIUserInfo struct {
	UserName  string `json:"userName" validate:"required"`
}

type APIUserPublic struct {
	UserName  string `json:"userName" validate:"required"`
	PublicKey string `json:"publicKey" validate:"required"`
}

type APIUserDetails struct {
	UserName   string `json:"userName" validate:"required"`
	PrivateKey string `json:"privateKey" validate:"required"`
	PublicKey  string `json:"publicKey" validate:"required"`
}

