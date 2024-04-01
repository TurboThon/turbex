package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	UserName   string             `json:"userName"`
	FirstName  string             `json:"firstName"`
	LastName   string             `json:"lastName"`
	Password   string             `json:"password"`
	PrivateKey string             `json:"privateKey"`
	PublicKey  string             `json:"publicKey"`
}

type APICreateUserRequest struct {
	UserName   string `json:"userName" validate:"required"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	Password   string `json:"password" validate:"required"`
	PrivateKey string `json:"privateKey" validate:"required"`
	PublicKey  string `json:"publicKey" validate:"required"`
}

type APIChangeUserRequest struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Password   string `json:"password"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type APIUserInfo struct {
	UserName  string `json:"userName" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type APIUserPublic struct {
	UserName  string `json:"userName" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	PublicKey string `json:"publicKey" validate:"required"`
}

type APIUserDetails struct {
	UserName   string `json:"userName" validate:"required"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	PrivateKey string `json:"privateKey" validate:"required"`
	PublicKey  string `json:"publicKey" validate:"required"`
}

type APIModifyUserRequest struct {
	UserName   string `json:"userName"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Password   string `json:"password"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}
