package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
  Id primitive.ObjectID `json:"id,omitempty"`
  FirstName string `json:"firstName" validate:"required"`
  LastName string `json:"lastName" validate:"required"`
  Password string `json:"password"`
  PrivateKey string `json:"privateKey"`
}
