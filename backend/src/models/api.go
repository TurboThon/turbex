package models

type APISuccess[D interface{}] struct {
  Data D `json:"data"`
}

type APIError struct {
  Error string `json:"error"`
}

